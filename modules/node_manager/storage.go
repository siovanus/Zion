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

package node_manager

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contract"
	"github.com/ethereum/go-ethereum/contract/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rlp"
)

var StartEpochID = common.Big1 // epoch started from 1, NOT 0!

var ErrEof = errors.New("EOF")

// storage key prefix
const (
	SKP_GLOBAL_CONFIG                 = "st_global_config"
	SKP_VALIDATOR                     = "st_validator"
	SKP_SIGNER_ADDR                   = "st_signer_addr"
	SKP_PROPOSAL_ADDR                 = "st_proposal_addr"
	SKP_ALL_VALIDATOR                 = "st_all_validator"
	SKP_TOTAL_POOL                    = "st_lock_pool"
	SKP_STAKE_INFO                    = "st_stake_info"
	SKP_UNLOCK_INFO                   = "st_unlock_info"
	SKP_CURRENT_EPOCH                 = "st_current_epoch"
	SKP_EPOCH_INFO                    = "st_epoch_info"
	SKP_ACCUMULATED_COMMISSION        = "st_accumulated_commission"
	SKP_VALIDATOR_ACCUMULATED_REWARDS = "st_validator_accumulated_rewards"
	SKP_VALIDATOR_OUTSTANDING_REWARDS = "st_validator_outstanding_rewards"
	SKP_OUTSTANDING_REWARDS           = "st_outstanding_rewards"
	SKP_VALIDATOR_SNAPSHOT_REWARDS    = "st_validator_snapshot_rewards"
	SKP_STAKE_STARTING_INFO           = "st_stake_starting_info"
	SKP_SIGN                          = "st_sign"
	SKP_SIGNER                        = "st_signer"
	SKP_COMMUNITY_INFO                = "st_community_info"
)

func setAccumulatedCommission(s *contract.ModuleContract, consensusAddr common.Address, accumulatedCommission *AccumulatedCommission) error {
	key := accumulatedCommissionKey(consensusAddr)
	store, err := rlp.EncodeToBytes(accumulatedCommission)
	if err != nil {
		return fmt.Errorf("setAccumulatedCommission, serialize accumulatedCommission error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func getAccumulatedCommission(s *contract.ModuleContract, consensusAddr common.Address) (*AccumulatedCommission, error) {
	accumulatedCommission := &AccumulatedCommission{}
	key := accumulatedCommissionKey(consensusAddr)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("getAccumulatedCommission, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, accumulatedCommission); err != nil {
		return nil, fmt.Errorf("getAccumulatedCommission, deserialize accumulatedCommission error: %v", err)
	}
	return accumulatedCommission, nil
}

func delAccumulatedCommission(s *contract.ModuleContract, consensusAddr common.Address) error {
	key := accumulatedCommissionKey(consensusAddr)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func setValidatorAccumulatedRewards(s *contract.ModuleContract, consensusAddr common.Address, validatorAccumulatedRewards *ValidatorAccumulatedRewards) error {
	key := validatorAccumulatedRewardsKey(consensusAddr)
	store, err := rlp.EncodeToBytes(validatorAccumulatedRewards)
	if err != nil {
		return fmt.Errorf("setValidatorAccumulatedRewards, serialize validatorAccumulatedRewards error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func getValidatorAccumulatedRewards(s *contract.ModuleContract, consensusAddr common.Address) (*ValidatorAccumulatedRewards, error) {
	validatorAccumulatedRewards := &ValidatorAccumulatedRewards{}
	key := validatorAccumulatedRewardsKey(consensusAddr)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetValidatorAccumulatedRewards, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, validatorAccumulatedRewards); err != nil {
		return nil, fmt.Errorf("GetValidatorAccumulatedRewards, deserialize validatorAccumulatedRewards error: %v", err)
	}
	return validatorAccumulatedRewards, nil
}

func delValidatorAccumulatedRewards(s *contract.ModuleContract, consensusAddr common.Address) error {
	key := validatorAccumulatedRewardsKey(consensusAddr)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func setValidatorOutstandingRewards(s *contract.ModuleContract, consensusAddr common.Address, validatorOutstandingRewards *ValidatorOutstandingRewards) error {
	key := validatorOutstandingRewardsKey(consensusAddr)
	store, err := rlp.EncodeToBytes(validatorOutstandingRewards)
	if err != nil {
		return fmt.Errorf("setValidatorOutstandingRewards, serialize validatorOutstandingRewards error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func getValidatorOutstandingRewards(s *contract.ModuleContract, consensusAddr common.Address) (*ValidatorOutstandingRewards, error) {
	validatorOutstandingRewards := &ValidatorOutstandingRewards{}
	key := validatorOutstandingRewardsKey(consensusAddr)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("getValidatorOutstandingRewards, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, validatorOutstandingRewards); err != nil {
		return nil, fmt.Errorf("getValidatorOutstandingRewards, deserialize validatorOutstandingRewards error: %v", err)
	}
	return validatorOutstandingRewards, nil
}

func delValidatorOutstandingRewards(s *contract.ModuleContract, consensusAddr common.Address) error {
	key := validatorOutstandingRewardsKey(consensusAddr)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func setOutstandingRewards(s *contract.ModuleContract, outstandingRewards *OutstandingRewards) error {
	key := outstandingRewardsKey()
	store, err := rlp.EncodeToBytes(outstandingRewards)
	if err != nil {
		return fmt.Errorf("setOutstandingRewards, serialize outstandingRewards error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func getOutstandingRewards(s *contract.ModuleContract) (*OutstandingRewards, error) {
	outstandingRewards := &OutstandingRewards{
		Rewards: NewDecFromBigInt(new(big.Int)),
	}
	key := outstandingRewardsKey()
	store, err := get(s, key)
	if err == ErrEof {
		return outstandingRewards, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getOutstandingRewards, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, outstandingRewards); err != nil {
		return nil, fmt.Errorf("getOutstandingRewards, deserialize outstandingRewards error: %v", err)
	}
	return outstandingRewards, nil
}

func increaseReferenceCount(s *contract.ModuleContract, consensusAddr common.Address, period uint64) error {
	validatorSnapshotRewards, err := getValidatorSnapshotRewards(s, consensusAddr, period)
	if err != nil {
		return fmt.Errorf("increaseReferenceCount, getValidatorSnapshotRewards error: %v", err)
	}
	if validatorSnapshotRewards.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	validatorSnapshotRewards.ReferenceCount++
	err = setValidatorSnapshotRewards(s, consensusAddr, period, validatorSnapshotRewards)
	if err != nil {
		return fmt.Errorf("increaseReferenceCount, setValidatorSnapshotRewards error: %v", err)
	}
	return nil
}

func decreaseReferenceCount(s *contract.ModuleContract, consensusAddr common.Address, period uint64) error {
	validatorSnapshotRewards, err := getValidatorSnapshotRewards(s, consensusAddr, period)
	if err != nil {
		return fmt.Errorf("decreaseReferenceCount, getValidatorSnapshotRewards error: %v", err)
	}
	if validatorSnapshotRewards.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	validatorSnapshotRewards.ReferenceCount--
	if validatorSnapshotRewards.ReferenceCount == 0 {
		delValidatorSnapshotRewards(s, consensusAddr, period)
	} else {
		err = setValidatorSnapshotRewards(s, consensusAddr, period, validatorSnapshotRewards)
		if err != nil {
			return fmt.Errorf("decreaseReferenceCount, setValidatorSnapshotRewards error: %v", err)
		}
	}
	return nil
}

func setValidatorSnapshotRewards(s *contract.ModuleContract, consensusAddr common.Address, period uint64, validatorSnapshotRewards *ValidatorSnapshotRewards) error {
	key := validatorSnapshotRewardsKey(consensusAddr, period)
	store, err := rlp.EncodeToBytes(validatorSnapshotRewards)
	if err != nil {
		return fmt.Errorf("setValidatorSnapshotRewards, serialize validatorSnapshotRewards error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func getValidatorSnapshotRewards(s *contract.ModuleContract, consensusAddr common.Address, period uint64) (*ValidatorSnapshotRewards, error) {
	validatorSnapshotRewards := &ValidatorSnapshotRewards{}
	key := validatorSnapshotRewardsKey(consensusAddr, period)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("getValidatorSnapshotRewards, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, validatorSnapshotRewards); err != nil {
		return nil, fmt.Errorf("getValidatorSnapshotRewards, deserialize validatorSnapshotRewards error: %v", err)
	}
	return validatorSnapshotRewards, nil
}

func delValidatorSnapshotRewards(s *contract.ModuleContract, consensusAddr common.Address, period uint64) error {
	key := validatorSnapshotRewardsKey(consensusAddr, period)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func setStakeStartingInfo(s *contract.ModuleContract, stakeAddress common.Address, consensusAddr common.Address, stakeStartingInfo *StakeStartingInfo) error {
	key := stakeStartingInfoKey(stakeAddress, consensusAddr)
	store, err := rlp.EncodeToBytes(stakeStartingInfo)
	if err != nil {
		return fmt.Errorf("setStakeStartingInfo, serialize stakeStartingInfo error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func getStakeStartingInfo(s *contract.ModuleContract, stakeAddress common.Address, consensusAddr common.Address) (*StakeStartingInfo, error) {
	stakeStartingInfo := &StakeStartingInfo{}
	key := stakeStartingInfoKey(stakeAddress, consensusAddr)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("getStakeStartingInfo, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, stakeStartingInfo); err != nil {
		return nil, fmt.Errorf("getStakeStartingInfo, deserialize stakeStartingInfo error: %v", err)
	}
	return stakeStartingInfo, nil
}

func delStakeStartingInfo(s *contract.ModuleContract, stakeAddress common.Address, consensusAddr common.Address) error {
	key := stakeStartingInfoKey(stakeAddress, consensusAddr)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func SetGlobalConfig(s *contract.ModuleContract, globalConfig *GlobalConfig) error {
	if globalConfig.MaxCommissionChange.Cmp(PercentDecimal) > 0 {
		return fmt.Errorf("SetGlobalConfig, MaxCommissionChange over size")
	}
	key := globalConfigKey()
	store, err := rlp.EncodeToBytes(globalConfig)
	if err != nil {
		return fmt.Errorf("setGlobalConfig, serialize globalConfig error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func GetGlobalConfigImpl(s *contract.ModuleContract) (*GlobalConfig, error) {
	globalConfig := &GlobalConfig{
		MaxCommissionChange:   GenesisMaxCommissionChange,
		MinInitialStake:       GenesisMinInitialStake,
		MinProposalStake:      GenesisMinProposalStake,
		BlockPerEpoch:         GenesisBlockPerEpoch,
		ConsensusValidatorNum: GenesisConsensusValidatorNum,
		VoterValidatorNum:     GenesisVoterValidatorNum,
	}
	key := globalConfigKey()
	store, err := get(s, key)
	if err == ErrEof {
		return globalConfig, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfigImpl, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, globalConfig); err != nil {
		return nil, fmt.Errorf("GetGlobalConfigImpl, deserialize globalConfig error: %v", err)
	}
	return globalConfig, nil
}

func addToAllValidators(s *contract.ModuleContract, consensusAddr common.Address) error {
	allValidators, err := getAllValidators(s)
	if err != nil {
		return fmt.Errorf("addToAllValidators, getAllValidators error: %v", err)
	}
	allValidators.AllValidators = append(allValidators.AllValidators, consensusAddr)
	if len(allValidators.AllValidators) > MaxValidatorNum {
		return fmt.Errorf("addToAllValidators, validator num is more than max")
	}
	err = setAllValidators(s, allValidators)
	if err != nil {
		return fmt.Errorf("addToAllValidators, set all validators error: %v", err)
	}
	return nil
}

func removeFromAllValidators(s *contract.ModuleContract, consensusAddr common.Address) error {
	allValidators, err := getAllValidators(s)
	if err != nil {
		return fmt.Errorf("removeFromAllValidators, getAllValidators error: %v", err)
	}
	j := 0
	for _, validator := range allValidators.AllValidators {
		if validator != consensusAddr {
			allValidators.AllValidators[j] = validator
			j++
		}
	}
	allValidators.AllValidators = allValidators.AllValidators[:j]
	err = setAllValidators(s, allValidators)
	if err != nil {
		return fmt.Errorf("removeFromAllValidators, set all validators error: %v", err)
	}
	return nil
}

func setValidator(s *contract.ModuleContract, validator *Validator) error {
	key := validatorKey(validator.ConsensusAddress)
	store, err := rlp.EncodeToBytes(validator)
	if err != nil {
		return fmt.Errorf("setValidator, serialize validator error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func delValidator(s *contract.ModuleContract, consensusAddr common.Address) error {
	key := validatorKey(consensusAddr)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func getValidator(s *contract.ModuleContract, consensusAddr common.Address) (*Validator, bool, error) {
	key := validatorKey(consensusAddr)
	store, err := get(s, key)
	if err == ErrEof {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("getValidator, get store error: %v", err)
	}
	validator := new(Validator)
	if err := rlp.DecodeBytes(store, validator); err != nil {
		return nil, false, fmt.Errorf("getValidator, deserialize validator error: %v", err)
	}
	return validator, true, nil
}

func setSignerAddr(s *contract.ModuleContract, signerAddr common.Address) error {
	key := signerAddrKey(signerAddr)
	_, err := get(s, key)
	if err != ErrEof {
		return fmt.Errorf("signer address aleady exist")
	}
	err = set(s, key, []byte{0x01})
	if err != nil {
		return err
	}
	return nil
}

func delSignerAddr(s *contract.ModuleContract, signerAddr common.Address) error {
	key := signerAddrKey(signerAddr)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func setProposalAddr(s *contract.ModuleContract, proposalAddr common.Address) error {
	key := proposalAddrKey(proposalAddr)
	_, err := get(s, key)
	if err != ErrEof {
		return fmt.Errorf("proposal address aleady exist")
	}
	err = set(s, key, []byte{0x01})
	if err != nil {
		return err
	}
	return nil
}

func delProposalAddr(s *contract.ModuleContract, proposalAddr common.Address) error {
	key := proposalAddrKey(proposalAddr)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func setAllValidators(s *contract.ModuleContract, allValidators *AllValidators) error {
	key := allValidatorKey()
	store, err := rlp.EncodeToBytes(allValidators)
	if err != nil {
		return fmt.Errorf("setAllValidators, serialize all validators error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func getAllValidators(s *contract.ModuleContract) (*AllValidators, error) {
	allValidators := &AllValidators{
		AllValidators: make([]common.Address, 0),
	}
	key := allValidatorKey()
	store, err := get(s, key)
	if err == ErrEof {
		return allValidators, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getAllValidators, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, allValidators); err != nil {
		return nil, fmt.Errorf("getAllValidators, deserialize all validators error: %v", err)
	}
	return allValidators, nil
}

func depositTotalPool(s *contract.ModuleContract, amount Dec) error {
	totalPool, err := getTotalPool(s)
	if err != nil {
		return fmt.Errorf("depositTotalPool, get total pool error: %v", err)
	}
	totalPool.TotalPool, err = totalPool.TotalPool.Add(amount)
	if err != nil {
		return fmt.Errorf("depositTotalPool, totalPool.TotalPool.Add error: %v", err)
	}
	err = setTotalPool(s, totalPool)
	if err != nil {
		return fmt.Errorf("depositTotalPool, setTotalPool error: %v", err)
	}
	return nil
}

func withdrawTotalPool(s *contract.ModuleContract, amount Dec) error {
	totalPool, err := getTotalPool(s)
	if err != nil {
		return fmt.Errorf("withdrawTotalPool, get total pool error: %v", err)
	}
	totalPool.TotalPool, err = totalPool.TotalPool.Sub(amount)
	if err != nil {
		return fmt.Errorf("withdrawTotalPool, totalPool.Sub error: %v", err)
	}
	err = setTotalPool(s, totalPool)
	if err != nil {
		return fmt.Errorf("withdrawTotalPool, setTotalPool error: %v", err)
	}
	return nil
}

func setTotalPool(s *contract.ModuleContract, totalPool *TotalPool) error {
	key := totalPoolKey()
	store, err := rlp.EncodeToBytes(totalPool)
	if err != nil {
		return fmt.Errorf("setStakeInfo, serialize stake info error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func getTotalPool(s *contract.ModuleContract) (*TotalPool, error) {
	totalPool := &TotalPool{NewDecFromBigInt(new(big.Int))}
	key := totalPoolKey()
	store, err := get(s, key)
	if err == ErrEof {
		return totalPool, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getTotalPool, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, totalPool); err != nil {
		return nil, fmt.Errorf("getTotalPool, deserialize totalPool error: %v", err)
	}
	return totalPool, nil
}

func setStakeInfo(s *contract.ModuleContract, stakeInfo *StakeInfo) error {
	key := stakeInfoKey(stakeInfo.StakeAddress, stakeInfo.ConsensusAddr)
	store, err := rlp.EncodeToBytes(stakeInfo)
	if err != nil {
		return fmt.Errorf("setStakeInfo, serialize stake info error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func delStakeInfo(s *contract.ModuleContract, stakeAddress common.Address, consensusAddr common.Address) error {
	key := stakeInfoKey(stakeAddress, consensusAddr)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func getStakeInfo(s *contract.ModuleContract, stakeAddress common.Address, consensusAddr common.Address) (*StakeInfo, bool, error) {
	stakeInfo := &StakeInfo{
		StakeAddress:  stakeAddress,
		ConsensusAddr: consensusAddr,
		Amount:        NewDecFromBigInt(new(big.Int)),
	}
	key := stakeInfoKey(stakeAddress, consensusAddr)
	store, err := get(s, key)
	if err == ErrEof {
		return stakeInfo, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("getStakeInfo, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, stakeInfo); err != nil {
		return nil, false, fmt.Errorf("getStakeInfo, deserialize stakeInfo error: %v", err)
	}
	return stakeInfo, true, nil
}

func addUnlockingInfo(s *contract.ModuleContract, stakeAddress common.Address, unlockingStake *UnlockingStake) error {
	unlockingInfo, err := getUnlockingInfo(s, stakeAddress)
	if err != nil {
		return fmt.Errorf("addUnlockingInfo, GetUnlockingInfo error: %v", err)
	}
	unlockingInfo.UnlockingStake = append(unlockingInfo.UnlockingStake, unlockingStake)
	if len(unlockingInfo.UnlockingStake) > MaxUnlockingNum {
		return fmt.Errorf("addUnlockingInfo, unlocking info more than max")
	}
	err = setUnlockingInfo(s, unlockingInfo)
	if err != nil {
		return fmt.Errorf("addUnlockingInfo, setUnlockingInfo error: %v", err)
	}
	return nil
}

func filterExpiredUnlockingInfo(s *contract.ModuleContract, stakeAddress common.Address) (Dec, error) {
	height := s.ContractRef().BlockHeight()
	unlockingInfo, err := getUnlockingInfo(s, stakeAddress)
	if err != nil {
		return Dec{nil}, fmt.Errorf("filterExpiredUnlockingInfo, GetUnlockingInfo error: %v", err)
	}
	j := 0
	expiredSum := NewDecFromBigInt(new(big.Int))
	for _, unlockingStake := range unlockingInfo.UnlockingStake {
		if unlockingStake.CompleteHeight.Cmp(height) == 1 {
			unlockingInfo.UnlockingStake[j] = unlockingStake
			j++
		} else {
			expiredSum, err = expiredSum.Add(unlockingStake.Amount)
			if err != nil {
				return Dec{nil}, fmt.Errorf("filterExpiredUnlockingInfo, expiredSum.Add error: %v", err)
			}
		}
	}
	unlockingInfo.UnlockingStake = unlockingInfo.UnlockingStake[:j]
	if len(unlockingInfo.UnlockingStake) == 0 {
		err = delUnlockingInfo(s, stakeAddress)
		if err != nil {
			return Dec{nil}, err
		}
	} else {
		err = setUnlockingInfo(s, unlockingInfo)
		if err != nil {
			return Dec{nil}, fmt.Errorf("filterExpiredUnlockingInfo, setUnlockingInfo error: %v", err)
		}
	}
	return expiredSum, nil
}

func setUnlockingInfo(s *contract.ModuleContract, unlockingInfo *UnlockingInfo) error {
	key := unlockingInfoKey(unlockingInfo.StakeAddress)
	store, err := rlp.EncodeToBytes(unlockingInfo)
	if err != nil {
		return fmt.Errorf("setUnlockingInfo, serialize unlock info error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func delUnlockingInfo(s *contract.ModuleContract, stakeAddress common.Address) error {
	key := unlockingInfoKey(stakeAddress)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func getUnlockingInfo(s *contract.ModuleContract, stakeAddress common.Address) (*UnlockingInfo, error) {
	unlockingInfo := &UnlockingInfo{
		StakeAddress:   stakeAddress,
		UnlockingStake: make([]*UnlockingStake, 0),
	}
	key := unlockingInfoKey(stakeAddress)
	store, err := get(s, key)
	if err == ErrEof {
		return unlockingInfo, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetUnlockingInfo, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, unlockingInfo); err != nil {
		return nil, fmt.Errorf("GetUnlockingInfo, deserialize unlocking info error: %v", err)
	}
	return unlockingInfo, nil
}

func setCurrentEpoch(s *contract.ModuleContract, ID *big.Int) error {
	key := currentEpochKey()
	err := set(s, key, ID.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func getCurrentEpoch(s *contract.ModuleContract) (*big.Int, error) {
	key := currentEpochKey()
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpoch, get store error: %v", err)
	}
	return new(big.Int).SetBytes(store), nil
}

func setCurrentEpochInfo(s *contract.ModuleContract, epochInfo *EpochInfo) error {
	// set current epoch
	err := setCurrentEpoch(s, epochInfo.ID)
	if err != nil {
		return err
	}
	//set epoch info
	err = setEpochInfo(s, epochInfo)
	if err != nil {
		return fmt.Errorf("setCurrentEpochInfo, setEpochInfo error: %v", err)
	}
	return nil
}

func setGenesisEpochInfo(db *state.StateDB, epochInfo *EpochInfo) error {
	s := state.NewStore(db, state.NewInfiniteGasMeter(), state.DefaultGasConfig())
	// set current epoch
	key1 := currentEpochKey()
	customSet(s, key1, epochInfo.ID.Bytes())
	//set epoch info
	key2 := epochInfoKey(epochInfo.ID)
	store, err := rlp.EncodeToBytes(epochInfo)
	if err != nil {
		return fmt.Errorf("setGenesisEpochInfo, serialize epoch info error: %v", err)
	}
	err = customSet(s, key2, store)
	if err != nil {
		return err
	}
	return nil
}

func GetCurrentEpochInfoImpl(s *contract.ModuleContract) (*EpochInfo, error) {
	ID, err := getCurrentEpoch(s)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfoImpl, getCurrentEpochInfo error: %v", err)
	}
	epochInfo, err := getEpochInfo(s, ID)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfoImpl, getEpochInfo error: %v", err)
	}
	return epochInfo, nil
}

func GetCurrentEpochInfoFromDB(db *state.StateDB) (*EpochInfo, error) {
	s := state.NewStore(db, state.NewInfiniteGasMeter(), state.DefaultGasConfig())
	key := currentEpochKey()
	store, err := customGet(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfoFromDB, get key store error: %v", err)
	}
	ID := new(big.Int).SetBytes(store)

	epochInfo := new(EpochInfo)
	key = epochInfoKey(ID)
	store, err = customGet(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfoFromDB, get info store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, epochInfo); err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfoFromDB, deserialize epoch info error: %v", err)
	}
	return epochInfo, nil
}

func setEpochInfo(s *contract.ModuleContract, epochInfo *EpochInfo) error {
	key := epochInfoKey(epochInfo.ID)
	store, err := rlp.EncodeToBytes(epochInfo)
	if err != nil {
		return fmt.Errorf("setEpochInfo, serialize epoch info error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func getEpochInfo(s *contract.ModuleContract, ID *big.Int) (*EpochInfo, error) {
	epochInfo := new(EpochInfo)
	key := epochInfoKey(ID)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetEpochInfo, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, epochInfo); err != nil {
		return nil, fmt.Errorf("GetEpochInfo, deserialize epoch info error: %v", err)
	}
	return epochInfo, nil
}

func SetCommunityInfo(s *contract.ModuleContract, communityInfo *CommunityInfo) error {
	if communityInfo.CommunityRate.Cmp(PercentDecimal) > 0 {
		return fmt.Errorf("setCommunityInfo, CommunityRate over size")
	}
	key := communityInfoKey()
	store, err := rlp.EncodeToBytes(communityInfo)
	if err != nil {
		return fmt.Errorf("setCommunityInfo, serialize community info error: %v", err)
	}
	err = set(s, key, store)
	if err != nil {
		return err
	}
	return nil
}

func GetCommunityInfoImpl(s *contract.ModuleContract) (*CommunityInfo, error) {
	communityInfo := &CommunityInfo{
		CommunityRate:    new(big.Int),
		CommunityAddress: common.EmptyAddress,
	}
	key := communityInfoKey()
	store, err := get(s, key)
	if err == ErrEof {
		return communityInfo, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetCommunityInfoImpl, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, communityInfo); err != nil {
		return nil, fmt.Errorf("GetCommunityInfoImpl, deserialize community info error: %v", err)
	}
	return communityInfo, nil
}

// ====================================================================
//
// `consensus sign` storage
//
// ====================================================================
func storeSign(s *contract.ModuleContract, sign *ConsensusSign) error {
	key := signKey(sign.Hash())
	value, err := rlp.EncodeToBytes(sign)
	if err != nil {
		return err
	}
	err = set(s, key, value)
	if err != nil {
		return err
	}
	return nil
}

func delSign(s *contract.ModuleContract, hash common.Hash) error {
	key := signKey(hash)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

func getSign(s *contract.ModuleContract, hash common.Hash) (*ConsensusSign, error) {
	key := signKey(hash)
	value, err := get(s, key)
	if err != nil {
		return nil, err
	}
	var sign *ConsensusSign
	if err := rlp.DecodeBytes(value, &sign); err != nil {
		return nil, err
	}
	return sign, nil
}

func storeSigner(s *contract.ModuleContract, hash common.Hash, signer common.Address) error {
	data, err := getSigners(s, hash)
	if err != nil {
		if err.Error() == ErrEof.Error() {
			data = make([]common.Address, 0)
		} else {
			return err
		}
	}
	data = append(data, signer)
	list := &AddressList{List: data}

	key := signerKey(hash)
	value, err := rlp.EncodeToBytes(list)
	if err != nil {
		return err
	}
	err = set(s, key, value)
	if err != nil {
		return err
	}

	return nil
}

func findSigner(s *contract.ModuleContract, hash common.Hash, signer common.Address) bool {
	list, err := getSigners(s, hash)
	if err != nil {
		return false
	}
	for _, v := range list {
		if v == signer {
			return true
		}
	}
	return false
}

func getSigners(s *contract.ModuleContract, hash common.Hash) ([]common.Address, error) {
	key := signerKey(hash)
	value, err := get(s, key)
	if err != nil {
		return nil, err
	}

	var list *AddressList
	if err := rlp.DecodeBytes(value, &list); err != nil {
		return nil, err
	}
	return list.List, nil
}

func getSignerSize(s *contract.ModuleContract, hash common.Hash) int {
	list, err := getSigners(s, hash)
	if err != nil {
		return 0
	}
	return len(list)
}

func clearSigner(s *contract.ModuleContract, hash common.Hash) error {
	key := signerKey(hash)
	err := del(s, key)
	if err != nil {
		return err
	}
	return nil
}

// ====================================================================
//
// storage basic operations
//
// ====================================================================

func get(s *contract.ModuleContract, key []byte) ([]byte, error) {
	return customGet(s.GetCacheDB(), key)
}

func set(s *contract.ModuleContract, key, value []byte) error {
	return customSet(s.GetCacheDB(), key, value)
}

func del(s *contract.ModuleContract, key []byte) error {
	return customDel(s.GetCacheDB(), key)
}

func customGet(s *state.Store, key []byte) ([]byte, error) {
	value, err := s.Get(key)
	if err != nil {
		return nil, err
	} else if value == nil || len(value) == 0 {
		return nil, ErrEof
	} else {
		return value, nil
	}
}

func customSet(s *state.Store, key, value []byte) error {
	err := s.Put(key, value)
	if err != nil {
		return err
	}
	return nil
}

func customDel(s *state.Store, key []byte) error {
	err := s.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

// ====================================================================
//
// storage keys
//
// ====================================================================

func globalConfigKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_GLOBAL_CONFIG))
}

func validatorKey(consensusAddr common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR), consensusAddr[:])
}

func signerAddrKey(signerAddr common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_SIGNER_ADDR), signerAddr[:])
}

func proposalAddrKey(proposalAddr common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_PROPOSAL_ADDR), proposalAddr[:])
}

func allValidatorKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_ALL_VALIDATOR))
}

func totalPoolKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_TOTAL_POOL))
}

func stakeInfoKey(stakeAddress common.Address, consensusAddr common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_STAKE_INFO), stakeAddress[:], consensusAddr[:])
}

func unlockingInfoKey(stakeAddress common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_UNLOCK_INFO), stakeAddress[:])
}

func currentEpochKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_CURRENT_EPOCH))
}

func epochInfoKey(ID *big.Int) []byte {
	return utils.ConcatKey(this, []byte(SKP_EPOCH_INFO), ID.Bytes())
}

func accumulatedCommissionKey(consensusAddr common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_ACCUMULATED_COMMISSION), consensusAddr[:])
}

func validatorAccumulatedRewardsKey(consensusAddr common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR_ACCUMULATED_REWARDS), consensusAddr[:])
}

func validatorOutstandingRewardsKey(consensusAddr common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR_OUTSTANDING_REWARDS), consensusAddr[:])
}

func outstandingRewardsKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_OUTSTANDING_REWARDS))
}

func validatorSnapshotRewardsKey(consensusAddr common.Address, period uint64) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR_SNAPSHOT_REWARDS), consensusAddr[:], utils.Uint64Bytes(period))
}

func stakeStartingInfoKey(stakeAddress common.Address, consensusAddr common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_STAKE_STARTING_INFO), stakeAddress[:], consensusAddr[:])
}

func signKey(hash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_SIGN), hash.Bytes())
}

func signerKey(hash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_SIGNER), hash.Bytes())
}

func communityInfoKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_COMMUNITY_INFO))
}
