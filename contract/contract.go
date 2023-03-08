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

package contract

import (
	"fmt"
	abiPkg "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

type (
	RegisterService func(native *ModuleContract)
	MethodHandler   func(contract *ModuleContract) ([]byte, error)
)

var (
	Contracts = make(map[common.Address]RegisterService)
)

// the gasUsage for the native contract transaction calculated according to the following formula:
// *		`gasUsage = gasRatio * gasTable[methodId]`
// the value in gas table for native tx is the max num for bench test in linux.
const (
	basicGas = uint64(21000) // minimum gas spent by transaction which failed before contract.handler, the default value is 21000 wei.
	gasRatio = float64(1.0)  // gasRatio is used to adjust the final value of gasUsage.
)

type ModuleContract struct {
	ref      *ContractRef
	db       *state.StateDB
	handlers map[string]MethodHandler // map method id to method handler
	gasTable map[string]uint64        // map method id to gas usage
	ab       *abiPkg.ABI
}

func NewModuleContract(db *state.StateDB, ref *ContractRef) *ModuleContract {
	return &ModuleContract{
		db:       db,
		ref:      ref,
		handlers: make(map[string]MethodHandler),
	}
}

func (s *ModuleContract) ContractRef() *ContractRef {
	return s.ref
}

func (s *ModuleContract) GetCacheDB() *state.CacheDB {
	return (*state.CacheDB)(s.db)
}

func (s *ModuleContract) StateDB() *state.StateDB {
	return s.db
}

func (s *ModuleContract) Prepare(ab *abiPkg.ABI, gasTb map[string]uint64) {
	s.ab = ab
	s.gasTable = make(map[string]uint64)
	for name, gas := range gasTb {
		id := MethodID(s.ab, name)
		final := uint64(float64(basicGas) + float64(gas)*gasRatio)
		s.gasTable[id] = final
	}
}

func (s *ModuleContract) Register(name string, handler MethodHandler) {
	methodID := MethodID(s.ab, name)
	s.handlers[methodID] = handler
}

// Invoke return execute ret and cost gas
func (s *ModuleContract) Invoke() ([]byte, error) {

	// pre-cost for failed tx which failed before `handler` execution.
	if gasLeft := s.ref.gasLeft; gasLeft < basicGas {
		s.ref.gasLeft = 0
		return nil, fmt.Errorf("gasLeft not enough, need %d, got %d", basicGas, gasLeft)
	} else {
		s.ref.gasLeft -= basicGas
	}

	// check context
	if !s.ref.CheckContexts() {
		return nil, fmt.Errorf("context error")
	}
	ctx := s.ref.CurrentContext()

	// find methodID
	if len(ctx.Payload) < 4 {
		return nil, fmt.Errorf("invalid input")
	}
	methodID := hexutil.Encode(ctx.Payload[:4])

	// register methods
	registerHandler, ok := Contracts[ctx.ContractAddress]
	if !ok {
		return nil, fmt.Errorf("failed to find contract: [%x]", ctx.ContractAddress)
	}
	registerHandler(s)

	// get method handler
	handler, ok := s.handlers[methodID]
	if !ok {
		return nil, fmt.Errorf("failed to find method: [%s]", methodID)
	}

	// check gas usage, the min value should be `basicGas`
	gasUsage, ok := s.gasTable[methodID]
	if !ok {
		return nil, fmt.Errorf("failed to find method: [%s]", methodID)
	}
	if gasUsage < basicGas {
		gasUsage = basicGas
	}
	// refund basic gas before tx get into `handler`
	s.ref.gasLeft += basicGas
	if gasLeft := s.ref.gasLeft; gasLeft < gasUsage {
		return nil, fmt.Errorf("gasLeft not enough, need %d, got %d", gasUsage, gasLeft)
	}

	// execute transaction and cost gas
	ret, err := handler(s)
	s.ref.gasLeft -= gasUsage
	return ret, err
}

func (s *ModuleContract) AddNotify(abi *abiPkg.ABI, topics []string, data ...interface{}) error {
	var topicIDs []common.Hash

	if topics == nil || len(topics) == 0 {
		return fmt.Errorf("AddNotify, topics length invalid")
	}

	topic := topics[0]
	topic, event, err := getTopicAndEvent(abi, topic)
	if err != nil {
		return fmt.Errorf("AddNotify, getTopicAndEvent err: %v", err)
	}
	topicIDs = append(topicIDs, event.ID)

	if len(data) != len(event.Inputs) {
		return fmt.Errorf("AddNotify, args length not equal to params number")
	}

	for i, input := range event.Inputs {
		if !input.Indexed {
			continue
		}

		topicID, ok := data[i].(common.Hash)
		if !ok {
			return fmt.Errorf("AddNotify, indexed field should be type of common.Hash")
		}
		topicIDs = append(topicIDs, topicID)
	}

	packedData, err := PackEvents(abi, topic, data...)
	if err != nil {
		return fmt.Errorf("AddNotify, PackEvents error: %v", err)
	}
	emitter := NewEventEmitter(s.ref.CurrentContext().ContractAddress, s.ContractRef().BlockHeight().Uint64(), s.StateDB())
	emitter.Event(topicIDs, packedData)

	return nil
}

func topic2CamelCase(topic string) string {
	return "evt" + abiPkg.ToCamelCase(topic)
}

func getTopicAndEvent(abi *abiPkg.ABI, topic string) (string, *abiPkg.Event, error) {
	eventInfo, ok := abi.Events[topic]
	if ok {
		return topic, &eventInfo, nil
	}

	topicWithCamel := topic2CamelCase(topic)
	eventInfo, ok = abi.Events[topicWithCamel]
	if ok {
		return topicWithCamel, &eventInfo, nil
	}
	return topic, nil, fmt.Errorf("topic %s not exist", topic)
}

// support module functions to evm functions.
type EVMHandler func(caller, addr common.Address, gas uint64, input []byte) ([]byte, uint64, error)

type ContractRef struct {
	contexts []*Context

	stateDB     *state.StateDB
	blockHeight *big.Int
	origin      common.Address
	txHash      common.Hash
	caller     common.Address
	evmHandler EVMHandler
	gasLeft    uint64
	value       *big.Int
	txTo        common.Address
}

func NewContractRef(
	db *state.StateDB,
	origin common.Address,
	caller common.Address,
	blockHeight *big.Int,
	txHash common.Hash,
	suppliedGas uint64,
	evmHandler EVMHandler) *ContractRef {

	return &ContractRef{
		contexts:    make([]*Context, 0),
		stateDB:     db,
		origin:      origin,
		caller:      caller,
		blockHeight: blockHeight,
		txHash:      txHash,
		gasLeft:     suppliedGas,
		evmHandler:  evmHandler,
		txTo:        common.EmptyAddress,
		value:       common.Big0,
	}
}

func (s *ContractRef) ModuleCall(
	caller,
	contractAddr common.Address,
	payload []byte,
) (ret []byte, gasLeft uint64, err error) {

	s.PushContext(&Context{
		Caller:          caller,
		ContractAddress: contractAddr,
		Payload:         payload,
	})
	defer s.PopContext()

	contract := NewModuleContract(s.stateDB, s)
	ret, err = contract.Invoke()
	gasLeft = s.gasLeft
	if err != nil {
		log.Error("Module contract", "invoke err", err, "txhash", s.txHash.Hex())
	}
	return
}

func (s *ContractRef) EVMCall(caller, contractAddr common.Address, gas uint64, input []byte) ([]byte, uint64, error) {
	if s.evmHandler == nil {
		return nil, 0, nil
	}
	return s.evmHandler(caller, contractAddr, gas, input)
}

func (s *ContractRef) SetValue(value *big.Int) {
	if value != nil && value.Cmp(common.Big0) > 0 {
		s.value = value
	}
}

// Value retrieve tx.value
func (s *ContractRef) Value() *big.Int {
	return s.value
}

func (s *ContractRef) SetTo(to common.Address) {
	if to != common.EmptyAddress {
		s.txTo = to
	}
}

// To retrieve tx.to
func (s *ContractRef) TxTo() common.Address {
	return s.txTo
}

func (s *ContractRef) StateDB() *state.StateDB {
	return s.stateDB
}

func (s *ContractRef) BlockHeight() *big.Int {
	return s.blockHeight
}

func (s *ContractRef) TxHash() common.Hash {
	return s.txHash
}

// MsgSender implement solidity grammar `msg.sender`
func (s *ContractRef) MsgSender() common.Address {
	return s.caller
}

// TxOrigin implement solidity grammar `tx.origin`
func (s *ContractRef) TxOrigin() common.Address {
	return s.origin
}

func (s *ContractRef) GasLeft() uint64 {
	return s.gasLeft
}

const (
	MAX_EXECUTE_CONTEXT = 128
)

type Context struct {
	Caller          common.Address
	ContractAddress common.Address
	Payload         []byte
}

// PushContext push current context to smart contract
func (s *ContractRef) PushContext(context *Context) {
	s.contexts = append(s.contexts, context)
}

// CurrentContext return smart contract current context
func (s *ContractRef) CurrentContext() *Context {
	if len(s.contexts) < 1 {
		return nil
	}
	return s.contexts[len(s.contexts)-1]
}

// PopContext pop smart contract current context
func (s *ContractRef) PopContext() {
	if len(s.contexts) > 1 {
		s.contexts = s.contexts[:len(s.contexts)-1]
	}
}

// CallingContext return smart contract caller context
func (s *ContractRef) CallingContext() *Context {
	if len(s.contexts) < 2 {
		return nil
	}
	return s.contexts[len(s.contexts)-2]
}

// EntryContext return smart contract entry entrance context
func (s *ContractRef) EntryContext() *Context {
	if len(s.contexts) < 1 {
		return nil
	}
	return s.contexts[0]
}

func (s *ContractRef) CheckContexts() bool {
	if len(s.contexts) == 0 {
		return false
	}
	if len(s.contexts) > MAX_EXECUTE_CONTEXT {
		return false
	}
	return true
}
