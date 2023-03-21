/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package backend

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/modules"
	"github.com/ethereum/go-ethereum/modules/utils"
	"math"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
)

// ===========================     utility function        ==========================

// callmsg implements core.Message to allow passing it as a transaction simulator.
type callmsg struct {
	ethereum.CallMsg
}

func (m callmsg) From() common.Address { return m.CallMsg.From }
func (m callmsg) Nonce() uint64        { return 0 }
func (m callmsg) CheckNonce() bool     { return false }
func (m callmsg) To() *common.Address  { return m.CallMsg.To }
func (m callmsg) GasPrice() *big.Int   { return m.CallMsg.GasPrice }
func (m callmsg) Gas() uint64          { return m.CallMsg.Gas }
func (m callmsg) Value() *big.Int      { return m.CallMsg.Value }
func (m callmsg) Data() []byte         { return m.CallMsg.Data }

// chain context
type chainContext struct {
	Chain  consensus.ChainHeaderReader
	engine consensus.Engine
}

func (c chainContext) Engine() consensus.Engine {
	return c.engine
}

// GetHeader blockContext need this function
func (c chainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	return c.Chain.GetHeader(hash, number)
}

const (
	systemGas      = math.MaxUint64 / 2 // system tx will be executed in evm, and gas calculating is needed.
	systemGasPrice = int64(0)           // consensus txs do not need to participate in gas price bidding
)

// getSystemMessage assemble system calling fields
func (s *backend) getSystemMessage(toAddress common.Address, data []byte, value *big.Int) callmsg {
	return callmsg{
		ethereum.CallMsg{
			From:     utils.SystemTxSender,
			Gas:      systemGas,
			GasPrice: big.NewInt(systemGasPrice),
			Value:    value,
			To:       &toAddress,
			Data:     data,
		},
	}
}

// getSystemCaller use fixed systemCaller as contract caller, and tx hash is useless in contract call.
func (s *backend) getSystemCaller(state *state.StateDB, height *big.Int) *modules.ContractRef {
	caller := utils.SystemTxSender
	hash := common.EmptyHash
	return modules.NewContractRef(state, caller, caller, height, hash, systemGas, nil)
}

// applyTransaction execute transaction without miner worker, and only succeed tx will be packed in block.
func (s *backend) applyTransaction(
	chain consensus.ChainHeaderReader,
	msg callmsg,
	state *state.StateDB,
	header *types.Header,
	chainContext core.ChainContext,
	commonTxs *[]*types.Transaction, receipts *[]*types.Receipt,
	sysTxs *[]*types.Transaction, usedGas *uint64, mining bool,
) (err error) {

	// check msg sender
	if msg.From() != utils.SystemTxSender {
		return fmt.Errorf("system tx sender invalid")
	}

	nonce := state.GetNonce(msg.From())
	expectedTx := types.NewTransaction(nonce, *msg.To(), msg.Value(), msg.Gas(), msg.GasPrice(), msg.Data())
	signer := types.MakeSigner(chain.Config(), header.Number)

	// miner worker execute `finalizeAndAssemble` in which the param of `mining` is true,  it's denote
	// that this tx comes from miner, and `validator` send governance tx in the same nonce is forbidden.
	// the sender of system tx is an unusual address, let the miner `signTx` to keep the signature in tx.Data
	// which denote that the system tx is mined by some one validator. this tx and the `actual` tx which
	// others sync node received should be compared and ensure that they are extreme the same.
	if mining {
		expectedTx, err = s.signer.SignTx(expectedTx, signer)
		if err != nil {
			return err
		}
	} else {
		// system tx CAN'T be nil or empty
		if sysTxs == nil || len(*sysTxs) == 0 || (*sysTxs)[0] == nil {
			return fmt.Errorf("supposed to get a actual transaction, but get none")
		}

		// check tx hash
		actualTx := (*sysTxs)[0]
		if expectedHash := signer.Hash(expectedTx); !bytes.Equal(signer.Hash(actualTx).Bytes(), expectedHash.Bytes()) {
			return fmt.Errorf("expected tx hash %v, nonce %d, to %s, value %s, gas %d, gasPrice %s, data %s;"+
				"get tx hash %v, nonce %d, to %s, value %s, gas %d, gasPrice %s, data %s",
				expectedHash.String(),
				expectedTx.Nonce(),
				expectedTx.To().String(),
				expectedTx.Value().String(),
				expectedTx.Gas(),
				expectedTx.GasPrice().String(),
				hex.EncodeToString(expectedTx.Data()),
				actualTx.Hash().String(),
				actualTx.Nonce(),
				actualTx.To().String(),
				actualTx.Value().String(),
				actualTx.Gas(),
				actualTx.GasPrice().String(),
				hex.EncodeToString(actualTx.Data()),
			)
		}

		// tx signature can be recovered and the sender should be equal to block `coinbase`
		sender, err := signer.Sender(actualTx)
		if err != nil {
			return fmt.Errorf("recover system tx sender failed, err: %v", err)
		}
		if sender != header.Coinbase {
			return fmt.Errorf("supposed to miner %s but got %s", header.Coinbase.Hex(), sender.Hex())
		}

		// reset tx and shift system tx list to next
		expectedTx = actualTx
		*sysTxs = (*sysTxs)[1:]
	}

	// execute system tx and get the receipt
	state.Prepare(expectedTx.Hash(), common.Hash{}, len(*commonTxs))
	gasUsed, err := applyMessage(msg, state, header, chain.Config(), chainContext)
	if err != nil {
		return err
	}
	*commonTxs = append(*commonTxs, expectedTx)
	var root []byte
	if chain.Config().IsByzantium(header.Number) {
		state.Finalise(true)
	} else {
		root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number)).Bytes()
	}
	*usedGas += gasUsed
	receipt := types.NewReceipt(root, false, *usedGas)
	receipt.TxHash = expectedTx.Hash()
	receipt.GasUsed = gasUsed

	// set the receipt logs and create a bloom for filtering
	receipt.Logs = state.GetLogs(expectedTx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = state.BlockHash()
	receipt.BlockNumber = header.Number
	receipt.TransactionIndex = uint(state.TxIndex())
	*receipts = append(*receipts, receipt)
	state.SetNonce(msg.From(), nonce+1)
	return nil
}

// applyMessage
func applyMessage(
	msg callmsg,
	state *state.StateDB,
	header *types.Header,
	chainConfig *params.ChainConfig,
	chainContext core.ChainContext,
) (uint64, error) {
	// Create a new context to be used in the EVM environment
	context := core.NewEVMBlockContext(header, chainContext, nil)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, vm.TxContext{Origin: msg.From(), GasPrice: big.NewInt(0)}, state, chainConfig, vm.Config{})
	// Apply the transaction to the current state (included in the env)
	ret, returnGas, err := vmenv.Call(
		vm.AccountRef(msg.From()),
		*msg.To(),
		msg.Data(),
		msg.Gas(),
		msg.Value(),
	)
	if err != nil {
		log.Error("apply message failed", "msg", string(ret), "err", err)
	}
	return msg.Gas() - returnGas, err
}

func packBlock(state *state.StateDB, chain consensus.ChainHeaderReader,
	header *types.Header, txs []*types.Transaction, receipts []*types.Receipt) *types.Block {

	var (
		block *types.Block
		root  common.Hash
	)

	// perform root calculation and block reorganization at the same time which with a large number of memory copy.
	// and reset the header root after actions done.
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
		wg.Done()
	}()
	go func() {
		// the header uncle hash will be settle as EmptyUncleHash which as the same of `nilUncleHash`
		block = types.NewBlock(header, txs, nil, receipts, trie.NewStackTrie(nil))
		wg.Done()
	}()
	wg.Wait()

	block.SetRoot(root)
	return block
}

type systemTxContext struct {
	chain    consensus.ChainHeaderReader
	state    *state.StateDB
	header   *types.Header
	chainCtx core.ChainContext
	txs      *[]*types.Transaction
	sysTxs   *[]*types.Transaction
	receipts *[]*types.Receipt
	usedGas  *uint64
	mining   bool
}

func (s *backend) executeTransaction(ctx *systemTxContext, contract common.Address, payload []byte) error {
	msg := s.getSystemMessage(contract, payload, common.Big0)
	return s.applyTransaction(ctx.chain, msg, ctx.state, ctx.header, ctx.chainCtx, ctx.txs, ctx.receipts, ctx.sysTxs, ctx.usedGas, ctx.mining)
}

func NewDefaultValSet(list []common.Address) hotstuff.ValidatorSet {
	return validator.NewSet(list, hotstuff.RoundRobin)
}
