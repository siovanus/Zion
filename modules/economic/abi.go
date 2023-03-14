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

package economic

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contract"
	"github.com/ethereum/go-ethereum/modules/cfg"
	. "github.com/ethereum/go-ethereum/modules/go_abi/economic_abi"
)

const contractName = "economic"

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(IEconomicABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = cfg.EconomicContractAddress
)

type MethodContractNameInput struct{}

func (m *MethodContractNameInput) Encode() ([]byte, error) {
	return contract.PackMethod(ABI, MethodName)
}
func (m *MethodContractNameInput) Decode(payload []byte) error { return nil }

type MethodContractNameOutput struct {
	Name string
}

func (m *MethodContractNameOutput) Encode() ([]byte, error) {
	m.Name = contractName
	return contract.PackOutputs(ABI, MethodName, m.Name)
}
func (m *MethodContractNameOutput) Decode(payload []byte) error {
	return contract.UnpackOutputs(ABI, MethodName, m, payload)
}

type MethodTotalSupplyInput struct{}

func (m *MethodTotalSupplyInput) Encode() ([]byte, error) {
	return contract.PackMethod(ABI, MethodTotalSupply)
}
func (m *MethodTotalSupplyInput) Decode(payload []byte) error { return nil }
