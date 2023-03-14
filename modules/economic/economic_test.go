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
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contract"
	. "github.com/ethereum/go-ethereum/modules/go_abi/economic_abi"
	"github.com/ethereum/go-ethereum/modules/node_manager"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	InitABI()
	node_manager.InitNodeManager()
	InitEconomic()
	os.Exit(m.Run())
}

func TestName(t *testing.T) {
	name := MethodName
	expect := contractName

	payload, err := new(MethodContractNameInput).Encode()
	assert.NoError(t, err)

	raw, err := contract.TestModuleCall(t, this, name, payload, common.Big0)
	assert.NoError(t, err)
	var got string
	assert.NoError(t, contract.UnpackOutputs(ABI, name, &got, raw))
	assert.Equal(t, expect, got)
}

// TestTotalSupply use command as follow to test each cases, and the result contains coverage and output. and use
// the flag of -count=1 to avoid the affect of test cache.
// cmd:
// go test -v -count=1 -cover github.com/ethereum/go-ethereum/contracts/native/economic -run TestTotalSupply
func TestTotalSupply(t *testing.T) {

	// genesis supply should be 100,000,000 and total supply has no upper limit.
	testcases := []struct {
		height  int
		expect  *big.Int
		testABI bool
	}{
		{0, big.NewInt(100000000), false},
		{40, big.NewInt(100000040), false},
		{200000000, big.NewInt(300000000), true},
	}
	name := MethodTotalSupply

	for _, tc := range testcases {
		var supply *big.Int

		payload, _ := new(MethodTotalSupplyInput).Encode()
		raw, err := contract.TestModuleCall(t, this, name, payload, common.Big0, tc.height)
		assert.NoError(t, err)

		if tc.testABI {
			output, err := ABI.Unpack(name, raw)
			assert.NoError(t, err)
			supply = *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
		} else {
			assert.NoError(t, contract.UnpackOutputs(ABI, name, &supply, raw))
		}

		got := new(big.Int).Div(supply, params.ZNT1)
		assert.Equal(t, tc.expect, got)
	}
}

func TestTransfer(t *testing.T) {
	var (
		from   = common.HexToAddress("0x123")
		to     = common.HexToAddress("0x456")
		amount = params.ZNT1
	)

	state := contract.NewTestStateDB()
	state.AddBalance(from, amount)

	_, ctx := contract.GenerateTestContext(t, common.Big0, from, state)
	if state.GetBalance(from).Cmp(amount) < 0 {
		t.Error("balance not enough")
	}
	state.SubBalance(from, amount)
	state.AddBalance(to, amount)
	t.Logf("base method `transfer` function %d", ctx.BreakPoint())
}
