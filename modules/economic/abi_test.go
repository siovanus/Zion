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
	"testing"

	"github.com/ethereum/go-ethereum/contract"
	. "github.com/ethereum/go-ethereum/modules/go_abi/economic_abi"

	"github.com/stretchr/testify/assert"
)

func TestABIMethodContractName(t *testing.T) {
	enc, err := contract.PackOutputs(ABI, MethodName, contractName)
	assert.NoError(t, err)
	params := new(MethodContractNameOutput)
	assert.NoError(t, contract.UnpackOutputs(ABI, MethodName, params, enc))
	assert.Equal(t, contractName, params.Name)
}

func TestABIMethodTotalSupply(t *testing.T) {
	expect := new(MethodTotalSupplyInput)
	enc, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodTotalSupplyInput)
	assert.NoError(t, got.Decode(enc))

	assert.Equal(t, expect, got)
}
