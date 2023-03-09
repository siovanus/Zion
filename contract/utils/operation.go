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
package utils

import (
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contract"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"math/big"
)

func ConcatKey(contract common.Address, args ...[]byte) []byte {
	temp := contract[:]
	for _, arg := range args {
		temp = append(temp, arg...)
	}
	return temp
}

func GetUint32Bytes(num uint32) []byte {
	var p [4]byte
	binary.LittleEndian.PutUint32(p[:], num)
	return p[:]
}

func GetBytesUint32(b []byte) uint32 {
	if len(b) != 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(b[:])
}

func GetBytesUint64(b []byte) uint64 {
	if len(b) != 8 {
		return 0
	}
	return binary.LittleEndian.Uint64(b[:])
}

func GetUint64Bytes(num uint64) []byte {
	var p [8]byte
	binary.LittleEndian.PutUint64(p[:], num)
	return p[:]
}

func ValidateOrigin(m *contract.ModuleContract, address common.Address) error {
	if m.ContractRef().TxOrigin() != address {
		return fmt.Errorf("ValidateOrigin, authentication failed!")
	}
	return nil
}

func ValidateSender(m *contract.ModuleContract, address common.Address) error {
	if m.ContractRef().MsgSender() != address {
		return fmt.Errorf("ValidateSender, authentication failed!")
	}
	return nil
}

func ModuleTransfer(s *state.StateDB, from, to common.Address, amount *big.Int) error {
	if amount.Sign() == -1 {
		return fmt.Errorf("amount can not be negative")
	}
	if !core.CanTransfer(s, from, amount) {
		return fmt.Errorf("%s insufficient balance", from.Hex())
	}
	core.Transfer(s, from, to, amount)
	return nil
}