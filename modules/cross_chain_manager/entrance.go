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

package cross_chain_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contract"
	"github.com/ethereum/go-ethereum/contract/utils"
	"github.com/ethereum/go-ethereum/modules/cfg"
	"github.com/ethereum/go-ethereum/modules/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/modules/cross_chain_manager/eth_common"
	"github.com/ethereum/go-ethereum/modules/cross_chain_manager/no_proof"
	"github.com/ethereum/go-ethereum/modules/cross_chain_manager/ripple"
	"github.com/ethereum/go-ethereum/modules/node_manager"
	"github.com/ethereum/go-ethereum/modules/side_chain_manager"
)

const (
	BLACKED_CHAIN = "BlackedChain"
)

// the real gas usage of `importOutTransfer` and `replenish` are 3291750 and 727125.
// in order to reduce the cross-chain cost, set them to be 300000 and 100000.
var (
	this     = cfg.ModuleContractAddrMap[cfg.ModuleCrossChain]
	gasTable = map[string]uint64{
		common.MethodContractName:        21000,
		common.MethodImportOuterTransfer: 300000,
		common.MethodBlackChain:          149625,
		common.MethodWhiteChain:          152250,
		common.MethodCheckDone:           57750,
		common.MethodReplenish:           100000,
		common.MethodMultiSignRipple:     100000,
		common.MethodReconstructRippleTx: 300000,
	}
)

var (
	NO_PROOF_ROUTER   = uint64(1)
	ETH_COMMON_ROUTER = uint64(2)

	RIPPLE_ROUTER = uint64(6)
)

func InitCrossChainManager() {
	contract.Contracts.RegisterContract(this, RegisterCrossChainManagerContract)
}

func RegisterCrossChainManagerContract(s *contract.ModuleContract) {
	s.Prepare(common.ABI, gasTable)

	s.Register(common.MethodContractName, Name)
	s.Register(common.MethodImportOuterTransfer, ImportOuterTransfer)
	s.Register(common.MethodBlackChain, BlackChain)
	s.Register(common.MethodWhiteChain, WhiteChain)
	s.Register(common.MethodCheckDone, CheckDone)
	s.Register(common.MethodReplenish, Replenish)

	// ripple
	s.Register(common.MethodMultiSignRipple, MultiSignRipple)
	s.Register(common.MethodReconstructRippleTx, ReconstructRippleTx)
}

func GetChainHandler(router uint64) (common.ChainHandler, error) {
	switch router {
	case NO_PROOF_ROUTER:
		return no_proof.NewNoProofHandler(), nil
	case ETH_COMMON_ROUTER:
		return eth_common.NewHandler(), nil
	case RIPPLE_ROUTER:
		return ripple.NewRippleHandler(), nil
	default:
		return nil, fmt.Errorf("not a supported router:%d", router)
	}
}

func Name(s *contract.ModuleContract) ([]byte, error) {
	return contract.PackOutputs(common.ABI, common.MethodContractName, cfg.ModuleCrossChain)
}

func CheckDone(s *contract.ModuleContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &common.CheckDoneParam{}
	if err := contract.UnpackMethod(common.ABI, common.MethodCheckDone, params, ctx.Payload); err != nil {
		return nil, err
	}
	if len(params.CrossChainID) == 0 || len(params.CrossChainID) > 2000 {
		return nil, fmt.Errorf("invalid cross chain id length, min 1, max 2000, current %v", len(params.CrossChainID))
	}
	err := common.CheckDoneTx(s, params.CrossChainID, params.ChainID)
	if err != nil && err != common.ErrTxAlreadyImported {
		return nil, err
	}
	return contract.PackOutputs(common.ABI, common.MethodCheckDone, err == common.ErrTxAlreadyImported)
}

func ImportOuterTransfer(s *contract.ModuleContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &common.EntranceParam{}
	if err := contract.UnpackMethod(common.ABI, common.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	srcChainID := params.SourceChainID
	blacked, err := CheckIfChainBlacked(s, srcChainID)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, CheckIfChainBlacked err: %v", err)
	}
	if blacked {
		return nil, fmt.Errorf("ImportExTransfer, source chain is blacked")
	}

	srcChain, err := side_chain_manager.GetSideChainObject(s, srcChainID)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, side_chain_manager.GetSideChain err: %v", err)
	} else if srcChain == nil {
		return nil, fmt.Errorf("ImportExTransfer, side chain %d is not registered", srcChainID)
	}

	handler, err := GetChainHandler(srcChain.Router)
	if err != nil {
		return nil, err
	}
	if handler == nil {
		return nil, fmt.Errorf("ImportExTransfer, handler for side chain %d is not exist", srcChainID)
	}

	txParam, err := handler.MakeDepositProposal(s)
	if err != nil {
		return nil, err
	}

	if txParam == nil {
		return contract.PackOutputs(common.ABI, common.MethodImportOuterTransfer, true)
	}

	//check target chain
	dstChainID := txParam.ToChainID
	if blacked, err = CheckIfChainBlacked(s, dstChainID); err != nil {
		return nil, fmt.Errorf("ImportExTransfer, CheckIfChainBlacked error: %v", err)
	}
	if blacked {
		return nil, fmt.Errorf("ImportExTransfer, target chain is blacked")
	}

	dstChain, err := side_chain_manager.GetSideChainObject(s, dstChainID)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, side_chain_manager.GetSideChain error: %v", err)
	}
	if dstChain == nil {
		return nil, fmt.Errorf("ImportExTransfer, side chain %d is not registered", dstChainID)
	}

	if dstChain.Router == RIPPLE_ROUTER {
		err := ripple.NewRippleHandler().MakeTransaction(s, txParam, srcChainID)
		if err != nil {
			return nil, err
		}
		return contract.PackOutputs(common.ABI, common.MethodImportOuterTransfer, true)
	}

	//NOTE, you need to store the tx in this
	if err := common.MakeTransaction(s, txParam, srcChainID); err != nil {
		return nil, err
	}

	return contract.PackOutputs(common.ABI, common.MethodImportOuterTransfer, true)
}

func MultiSignRipple(s *contract.ModuleContract) ([]byte, error) {
	handler := ripple.NewRippleHandler()

	//1. multi sign
	err := handler.MultiSign(s)
	if err != nil {
		return nil, err
	}
	return contract.PackOutputs(common.ABI, common.MethodMultiSignRipple, true)
}

func ReconstructRippleTx(s *contract.ModuleContract) ([]byte, error) {
	handler := ripple.NewRippleHandler()

	err := handler.ReconstructTx(s)
	if err != nil {
		return nil, err
	}
	return contract.PackOutputs(common.ABI, common.MethodReconstructRippleTx, true)
}

func BlackChain(s *contract.ModuleContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &common.BlackChainParam{}
	if err := contract.UnpackMethod(common.ABI, common.MethodBlackChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(s, common.MethodBlackChain, utils.GetUint64Bytes(params.ChainID), s.ContractRef().MsgSender(), node_manager.Signer)
	if err != nil {
		return nil, fmt.Errorf("BlackChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return contract.PackOutputs(common.ABI, common.MethodBlackChain, true)
	}

	PutBlackChain(s, params.ChainID)
	return contract.PackOutputs(common.ABI, common.MethodBlackChain, true)
}

func WhiteChain(s *contract.ModuleContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &common.BlackChainParam{}
	if err := contract.UnpackMethod(common.ABI, common.MethodWhiteChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(s, common.MethodWhiteChain, ctx.Payload, s.ContractRef().MsgSender(), node_manager.Signer)
	if err != nil {
		return nil, fmt.Errorf("WhiteChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return contract.PackOutputs(common.ABI, common.MethodWhiteChain, true)
	}

	RemoveBlackChain(s, params.ChainID)
	return contract.PackOutputs(common.ABI, common.MethodWhiteChain, true)
}

func Replenish(s *contract.ModuleContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &common.ReplenishParam{}
	if err := contract.UnpackMethod(common.ABI, common.MethodReplenish, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("Replenish, unpack params error: %s", err)
	}

	if len(params.TxHashes) == 0 || len(params.TxHashes) > 200 {
		return nil, fmt.Errorf("invalid replenish hash length, min 1, max 200, current %v", len(params.TxHashes))
	}
	err := common.NotifyReplenish(s, params.TxHashes, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("Replenish, NotifyReplenish error: %s", err)
	}
	return contract.PackOutputs(common.ABI, common.MethodReplenish, true)
}
