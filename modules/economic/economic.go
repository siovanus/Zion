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
	"github.com/ethereum/go-ethereum/contract"
	"github.com/ethereum/go-ethereum/modules/cfg"
	. "github.com/ethereum/go-ethereum/modules/go_abi/economic_abi"
	"github.com/ethereum/go-ethereum/modules/node_manager"
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

var (
	gasTable = map[string]uint64{
		MethodName:        39375,
		MethodTotalSupply: 23625,
		MethodReward:      73500,
	}
)

var (
	RewardPerBlock = params.ZNT1
	GenesisSupply  = params.GenesisSupply
)

func InitEconomic() {
	InitABI()
	contract.Contracts[this] = RegisterEconomicContract
}

func RegisterEconomicContract(s *contract.ModuleContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodTotalSupply, TotalSupply)
	s.Register(MethodReward, Reward)
}

func Name(s *contract.ModuleContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

func TotalSupply(s *contract.ModuleContract) ([]byte, error) {
	height := s.ContractRef().BlockHeight()

	supply := GenesisSupply
	if height.Uint64() > 0 {
		reward := new(big.Int).Mul(height, RewardPerBlock)
		supply = new(big.Int).Add(supply, reward)
	}
	return contract.PackOutputs(ABI, MethodTotalSupply, supply)
}

func Reward(s *contract.ModuleContract) ([]byte, error) {

	community, err := node_manager.GetCommunityInfoFromDB(s.StateDB())
	if err != nil {
		return nil, fmt.Errorf("GetCommunityInfo failed, err: %v", err)
	}

	// allow empty address as reward pool
	poolAddr := community.CommunityAddress
	rewardPerBlock := node_manager.NewDecFromBigInt(RewardPerBlock)
	rewardFactor := node_manager.NewDecFromBigInt(community.CommunityRate)
	poolRwdAmt, err := rewardPerBlock.MulWithPercentDecimal(rewardFactor)
	if err != nil {
		return nil, fmt.Errorf("Calculate pool reward amount failed, err: %v ", err)
	}
	stakingRwdAmt, err := rewardPerBlock.Sub(poolRwdAmt)
	if err != nil {
		return nil, fmt.Errorf("Calculate staking reward amount, failed, err: %v ", err)
	}

	poolRwd := &RewardAmount{
		Address: poolAddr,
		Amount:  poolRwdAmt.BigInt(),
	}
	stakingRwd := &RewardAmount{
		Address: cfg.NodeManagerContractAddress,
		Amount:  stakingRwdAmt.BigInt(),
	}

	output := new(MethodRewardOutput)
	output.List = []*RewardAmount{poolRwd, stakingRwd}
	return output.Encode()
}