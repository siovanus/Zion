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

package testutils

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/tool"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

func GenesisBlock(g *core.Genesis, db ethdb.Database) *types.Block {
	return g.MustCommit(db)
}

func GenesisAndKeys(n int) (*core.Genesis, []*ecdsa.PrivateKey, error) {
	list := make([]*ecdsa.PrivateKey, n)
	vals := make([]common.Address, n)

	for i := 0; i < n; i++ {
		pk, err := crypto.GenerateKey()
		if err != nil {
			return nil, nil, err
		}
		list[i] = pk
		addr := crypto.PubkeyToAddress(pk.PublicKey)
		vals[i] = addr
	}

	g, err := Genesis(vals)
	if err != nil {
		return nil, nil, err
	}
	return g, list, nil
}

func Genesis(validators []common.Address) (*core.Genesis, error) {
	g := new(core.Genesis)
	g.Config = &params.ChainConfig{
		ChainID:  new(big.Int).SetUint64(params.MainnetChainID),
		HotStuff: &params.HotStuffConfig{Protocol: "base"},
	}
	g.Alloc = core.GenesisAlloc{
		validators[0]: core.GenesisAccount{
			Balance: params.GenesisSupply,
		},
	}
	g.Governance = []core.GovernanceAccount{}
	for _, v := range validators {
		g.Governance = append(g.Governance, core.GovernanceAccount{
			Validator: v,
		})
	}
	g.Difficulty = big.NewInt(1)
	extra, err := tool.EncodeGenesisExtra(validators)
	if err != nil {
		return nil, err
	}
	rawExtra, err := hexutil.Decode(extra)
	if err != nil {
		return nil, err
	}
	g.ExtraData = rawExtra
	g.GasLimit = 30000000
	g.Timestamp = 1638385012 // 2021
	return g, nil
}
