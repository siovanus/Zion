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

package proposal_manager

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/proposal_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/community"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	PROPOSE_EVENT            = "Propose"
	PROPOSE_CONFIG_EVENT     = "ProposeConfig"
	PROPOSE_COMMUNITY_EVENT  = "ProposeCommunity"
	PROPOSE_SIDE_CHAIN_EVENT = "ProposeSideChain"
	VOTE_PROPOSAL_EVENT      = "VoteProposal"

	MaxContentLength int = 4000
)

var (
	gasTable = map[string]uint64{
		MethodPropose:                  979125,
		MethodProposeConfig:            756000,
		MethodProposeCommunity:         693000,
		MethodProposeSideChain:         693000,
		MethodVoteProposal:             603750,
		MethodGetProposal:              118125,
		MethodGetProposalList:          94500,
		MethodGetConfigProposalList:    73500,
		MethodGetCommunityProposalList: 84000,
		MethodGetSideChainProposalList: 84000,
	}
)

func InitProposalManager() {
	InitABI()
	native.Contracts[this] = RegisterProposalManagerContract
}

func RegisterProposalManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodPropose, Propose)
	s.Register(MethodProposeConfig, ProposeConfig)
	s.Register(MethodProposeCommunity, ProposeCommunity)
	s.Register(MethodProposeSideChain, ProposeSideChain)
	s.Register(MethodVoteProposal, VoteProposal)
	s.Register(MethodGetProposal, GetProposal)
	s.Register(MethodGetProposalList, GetProposalList)
	s.Register(MethodGetConfigProposalList, GetConfigProposalList)
	s.Register(MethodGetCommunityProposalList, GetCommunityProposalList)
	s.Register(MethodGetSideChainProposalList, GetSideChainProposalList)
}

func Propose(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller
	value := s.ContractRef().Value()
	toAddress := s.ContractRef().TxTo()

	if ctx.Caller != s.ContractRef().TxOrigin() {
		return nil, fmt.Errorf("Propose, contract call forbidden")
	}
	globalConfig, err := node_manager.GetGlobalConfigImpl(s)
	if err != nil {
		return nil, fmt.Errorf("Propose, GetGlobalConfigImpl error: %v", err)
	}
	if toAddress != utils.ProposalManagerContractAddress {
		return nil, fmt.Errorf("Propose, to address %x must be proposal manager contract address %x", toAddress, utils.ProposalManagerContractAddress)
	}
	if value.Cmp(globalConfig.MinProposalStake) == -1 {
		return nil, fmt.Errorf("Propose, value is less than globalConfig.MinProposalStake")
	}

	params := &ProposeParam{}
	if err := utils.UnpackMethod(ABI, MethodPropose, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("Propose, unpack params error: %v", err)
	}

	if len(params.Content) > MaxContentLength {
		return nil, fmt.Errorf("Propose, content is more than max length")
	}

	// remove expired proposal
	err = removeExpiredFromProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("Propose, removeExpiredFromProposalList error: %v", err)
	}

	proposalID, err := getProposalID(s)
	if err != nil {
		return nil, fmt.Errorf("Propose, getProposalID error: %v", err)
	}
	proposalList, err := getProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("Propose, getProposalList error: %v", err)
	}
	if len(proposalList.ProposalList) >= ProposalListLen {
		return nil, fmt.Errorf("Propose, proposal is more than max length %d", ProposalListLen)
	}
	proposal := &Proposal{
		ID:        proposalID,
		Address:   ctx.Caller,
		Type:      Normal,
		Content:   params.Content,
		EndHeight: new(big.Int).Add(height, globalConfig.BlockPerEpoch),
		Stake:     value,
	}
	proposalList.ProposalList = append(proposalList.ProposalList, proposal.ID)
	err = setProposalList(s, proposalList)
	if err != nil {
		return nil, fmt.Errorf("Propose, setProposalList error: %v", err)
	}
	err = setProposal(s, proposal)
	if err != nil {
		return nil, fmt.Errorf("Propose, setProposal error: %v", err)
	}
	setProposalID(s, new(big.Int).Add(proposalID, common.Big1))

	err = s.AddNotify(ABI, []string{PROPOSE_EVENT}, proposal.ID.String(), caller.Hex(), proposal.Stake.String(), hex.EncodeToString(params.Content))
	if err != nil {
		return nil, fmt.Errorf("Propose, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodPropose, true)
}

func ProposeConfig(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller
	value := s.ContractRef().Value()
	toAddress := s.ContractRef().TxTo()

	if ctx.Caller != s.ContractRef().TxOrigin() {
		return nil, fmt.Errorf("ProposeConfig, contract call forbidden")
	}
	globalConfig, err := node_manager.GetGlobalConfigImpl(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeConfig, GetGlobalConfigImpl error: %v", err)
	}
	if toAddress != utils.ProposalManagerContractAddress {
		return nil, fmt.Errorf("ProposeConfig, to address %x must be proposal manager contract address %x", toAddress, utils.ProposalManagerContractAddress)
	}
	if value.Cmp(globalConfig.MinProposalStake) == -1 {
		return nil, fmt.Errorf("ProposeConfig, value is less than globalConfig.MinProposalStake")
	}

	params := &ProposeConfigParam{}
	if err := utils.UnpackMethod(ABI, MethodProposeConfig, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("ProposeConfig, unpack params error: %v", err)
	}

	if len(params.Content) > MaxContentLength {
		return nil, fmt.Errorf("ProposeConfig, content is more than max length")
	}

	config := new(node_manager.GlobalConfig)
	err = rlp.DecodeBytes(params.Content, config)
	if err != nil {
		return nil, fmt.Errorf("ProposeConfig, deserialize global config error: %v", err)
	}

	if config.ConsensusValidatorNum != 0 && config.ConsensusValidatorNum < node_manager.GenesisConsensusValidatorNum {
		return nil, fmt.Errorf("ProposeConfig, consensus num is less than %d", node_manager.GenesisConsensusValidatorNum)
	}
	if config.BlockPerEpoch.Cmp(node_manager.MinBlockPerEpoch) < 0 {
		return nil, fmt.Errorf("ProposeConfig, block per epoch is less than %d", node_manager.MinBlockPerEpoch)
	}
	if config.MaxCommissionChange.Cmp(node_manager.GenesisMaxCommissionChange) > 0 {
		return nil, fmt.Errorf("ProposeConfig, MaxCommissionChange is more than %d", node_manager.GenesisMaxCommissionChange)
	}
	if config.MinInitialStake.Sign() < 0 {
		return nil, fmt.Errorf("ProposeConfig, MinInitialStake is negative")
	}
	if config.MinProposalStake.Sign() < 0 {
		return nil, fmt.Errorf("ProposeConfig, MinProposalStake is negative")
	}

	// remove expired proposal
	err = removeExpiredFromConfigProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeConfig, removeExpiredFromConfigProposalList error: %v", err)
	}

	proposalID, err := getProposalID(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeConfig, getProposalID error: %v", err)
	}
	configProposalList, err := getConfigProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeConfig, getConfigProposalList error: %v", err)
	}
	if len(configProposalList.ConfigProposalList) >= ProposalListLen {
		return nil, fmt.Errorf("ProposeConfig, proposal is more than max length %d", ProposalListLen)
	}
	proposal := &Proposal{
		ID:        proposalID,
		Address:   ctx.Caller,
		Type:      UpdateGlobalConfig,
		Content:   params.Content,
		EndHeight: new(big.Int).Add(height, globalConfig.BlockPerEpoch),
		Stake:     value,
	}
	configProposalList.ConfigProposalList = append(configProposalList.ConfigProposalList, proposal.ID)
	err = setConfigProposalList(s, configProposalList)
	if err != nil {
		return nil, fmt.Errorf("ProposeConfig, setConfigProposalList error: %v", err)
	}
	err = setProposal(s, proposal)
	if err != nil {
		return nil, fmt.Errorf("ProposeConfig, setProposal error: %v", err)
	}
	setProposalID(s, new(big.Int).Add(proposalID, common.Big1))

	err = s.AddNotify(ABI, []string{PROPOSE_CONFIG_EVENT}, proposal.ID.String(), caller.Hex(), proposal.Stake.String(), hex.EncodeToString(params.Content))
	if err != nil {
		return nil, fmt.Errorf("ProposeConfig, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodProposeConfig, true)
}

func ProposeCommunity(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller
	value := s.ContractRef().Value()
	toAddress := s.ContractRef().TxTo()

	if ctx.Caller != s.ContractRef().TxOrigin() {
		return nil, fmt.Errorf("ProposeCommunity, contract call forbidden")
	}
	globalConfig, err := node_manager.GetGlobalConfigImpl(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeCommunity, GetGlobalConfigImpl error: %v", err)
	}
	if toAddress != utils.ProposalManagerContractAddress {
		return nil, fmt.Errorf("ProposeCommunity, to address %x must be proposal manager contract address %x", toAddress, utils.ProposalManagerContractAddress)
	}
	if value.Cmp(globalConfig.MinProposalStake) == -1 {
		return nil, fmt.Errorf("ProposeCommunity, value is less than globalConfig.MinProposalStake")
	}

	params := &ProposeCommunityParam{}
	if err := utils.UnpackMethod(ABI, MethodProposeCommunity, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("ProposeCommunity, unpack params error: %v", err)
	}

	if len(params.Content) > MaxContentLength {
		return nil, fmt.Errorf("ProposeCommunity, content is more than max length")
	}

	info := new(community.CommunityInfo)
	err = rlp.DecodeBytes(params.Content, info)
	if err != nil {
		return nil, fmt.Errorf("ProposeCommunity, deserialize community info error: %v", err)
	}
	if info.CommunityRate.Sign() == -1 {
		return nil, fmt.Errorf("ProposeCommunity, communityRate is negative")
	}
	if info.CommunityRate.Cmp(node_manager.PercentDecimal) == 1 {
		return nil, fmt.Errorf("ProposeCommunity, communityRate can not more than 100 percent")
	}

	// remove expired proposal
	err = removeExpiredFromCommunityProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeCommunity, removeExpiredFromCommunityProposalList error: %v", err)
	}

	proposalID, err := getProposalID(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeCommunity, getProposalID error: %v", err)
	}
	communityProposalList, err := getCommunityProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeCommunity, getCommunityProposalList error: %v", err)
	}
	if len(communityProposalList.CommunityProposalList) >= ProposalListLen {
		return nil, fmt.Errorf("ProposeCommunity, proposal is more than max length %d", ProposalListLen)
	}
	proposal := &Proposal{
		ID:        proposalID,
		Address:   ctx.Caller,
		Type:      UpdateCommunityInfo,
		Content:   params.Content,
		EndHeight: new(big.Int).Add(height, globalConfig.BlockPerEpoch),
		Stake:     value,
	}
	communityProposalList.CommunityProposalList = append(communityProposalList.CommunityProposalList, proposal.ID)
	err = setCommunityProposalList(s, communityProposalList)
	if err != nil {
		return nil, fmt.Errorf("ProposeCommunity, setCommunityProposalList error: %v", err)
	}
	err = setProposal(s, proposal)
	if err != nil {
		return nil, fmt.Errorf("ProposeCommunity, setProposal error: %v", err)
	}
	setProposalID(s, new(big.Int).Add(proposalID, common.Big1))

	err = s.AddNotify(ABI, []string{PROPOSE_COMMUNITY_EVENT}, proposal.ID.String(), caller.Hex(), proposal.Stake.String(), hex.EncodeToString(params.Content))
	if err != nil {
		return nil, fmt.Errorf("ProposeCommunity, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodProposeCommunity, true)
}

func ProposeSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller
	value := s.ContractRef().Value()
	toAddress := s.ContractRef().TxTo()

	if ctx.Caller != s.ContractRef().TxOrigin() {
		return nil, fmt.Errorf("ProposeSideChain, contract call forbidden")
	}
	globalConfig, err := node_manager.GetGlobalConfigImpl(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeSideChain, GetGlobalConfigImpl error: %v", err)
	}
	if toAddress != utils.ProposalManagerContractAddress {
		return nil, fmt.Errorf("ProposeSideChain, to address %x must be proposal manager contract address %x", toAddress, utils.ProposalManagerContractAddress)
	}
	if value.Cmp(globalConfig.MinProposalStake) == -1 {
		return nil, fmt.Errorf("ProposeSideChain, value is less than globalConfig.MinProposalStake")
	}

	params := &ProposeSideChainParam{}
	if err := utils.UnpackMethod(ABI, MethodProposeSideChain, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("ProposeSideChain, unpack params error: %v", err)
	}

	if len(params.Content) > MaxContentLength {
		return nil, fmt.Errorf("ProposeSideChain, content is more than max length")
	}

	sideChainInfo := new(side_chain_manager.SideChain)
	err = rlp.DecodeBytes(params.Content, sideChainInfo)
	if err != nil {
		return nil, fmt.Errorf("ProposeSideChain, deserialize side chain info error: %v", err)
	}
	if len(sideChainInfo.Name) > 100 {
		return nil, fmt.Errorf("param name too long, max is 100")
	}
	if len(sideChainInfo.ExtraInfo) > 1000000 {
		return nil, fmt.Errorf("param extra info too long, max is 1000000")
	}
	if len(sideChainInfo.CCMCAddress) > 1000 {
		return nil, fmt.Errorf("ccmc address info too long, max is 1000")
	}

	// remove expired proposal
	err = removeExpiredFromSideChainProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeSideChain, removeExpiredFromSideChainProposalList error: %v", err)
	}

	proposalID, err := getProposalID(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeSideChain, getProposalID error: %v", err)
	}
	sideChainProposalList, err := getSideChainProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("ProposeSideChain, getSideChainProposalList error: %v", err)
	}
	if len(sideChainProposalList.SideChainProposalList) >= ProposalListLen {
		return nil, fmt.Errorf("ProposeSideChain, proposal is more than max length %d", ProposalListLen)
	}
	proposal := &Proposal{
		ID:        proposalID,
		Address:   ctx.Caller,
		Type:      UpdateSideChain,
		Content:   params.Content,
		EndHeight: new(big.Int).Add(height, globalConfig.BlockPerEpoch),
		Stake:     value,
	}
	sideChainProposalList.SideChainProposalList = append(sideChainProposalList.SideChainProposalList, proposal.ID)
	err = setSideChainProposalList(s, sideChainProposalList)
	if err != nil {
		return nil, fmt.Errorf("ProposeSideChain, setSideChainProposalList error: %v", err)
	}
	err = setProposal(s, proposal)
	if err != nil {
		return nil, fmt.Errorf("ProposeSideChain, setProposal error: %v", err)
	}
	setProposalID(s, new(big.Int).Add(proposalID, common.Big1))

	err = s.AddNotify(ABI, []string{PROPOSE_SIDE_CHAIN_EVENT}, proposal.ID.String(), caller.Hex(), proposal.Stake.String(), hex.EncodeToString(params.Content))
	if err != nil {
		return nil, fmt.Errorf("ProposeSideChain, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodProposeSideChain, true)
}

func VoteProposal(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &VoteProposalParam{}
	if err := utils.UnpackMethod(ABI, MethodVoteProposal, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("VoteProposal, unpack params error: %v", err)
	}

	proposal, err := getProposal(s, params.ID)
	if err != nil {
		return nil, fmt.Errorf("VoteProposal, getProposal error: %v", err)
	}

	if proposal.Status == PASS {
		return utils.PackOutputs(ABI, MethodVoteProposal, true)
	}
	if proposal.Status == FAIL || proposal.EndHeight.Cmp(s.ContractRef().BlockHeight()) < 0 {
		return nil, fmt.Errorf("VoteProposal, proposal already failed")
	}

	success, err := node_manager.CheckConsensusSigns(s, MethodVoteProposal, ctx.Payload, caller, node_manager.Proposer)
	if err != nil {
		return nil, fmt.Errorf("VoteProposal, node_manager.CheckConsensusSigns error: %v", err)
	}
	if success {
		// update proposal status
		proposal.Status = PASS
		err = setProposal(s, proposal)
		if err != nil {
			return nil, fmt.Errorf("VoteProposal, setProposal error: %v", err)
		}

		// transfer token
		err = contract.NativeTransfer(s.StateDB(), this, proposal.Address, proposal.Stake)
		if err != nil {
			return nil, fmt.Errorf("Propose, utils.NativeTransfer error: %v", err)
		}

		communityInfo, err := community.GetCommunityInfoImpl(s)
		if err != nil {
			return nil, fmt.Errorf("VoteProposal, node_manager.GetCommunityInfoImpl error: %v", err)
		}

		switch proposal.Type {
		case UpdateGlobalConfig:
			config := new(node_manager.GlobalConfig)
			err := rlp.DecodeBytes(proposal.Content, config)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, deserialize global config error: %v", err)
			}

			globalConfig, err := node_manager.GetGlobalConfigImpl(s)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, node_manager.GetGlobalConfigImpl error: %v", err)
			}
			if config.ConsensusValidatorNum >= node_manager.GenesisConsensusValidatorNum {
				globalConfig.ConsensusValidatorNum = config.ConsensusValidatorNum
			}
			if config.VoterValidatorNum > 0 {
				globalConfig.VoterValidatorNum = config.VoterValidatorNum
			}
			if globalConfig.ConsensusValidatorNum < globalConfig.VoterValidatorNum {
				globalConfig.VoterValidatorNum = globalConfig.ConsensusValidatorNum
			}
			if config.BlockPerEpoch.Cmp(node_manager.MinBlockPerEpoch) > 0 {
				globalConfig.BlockPerEpoch = config.BlockPerEpoch
			}
			if config.MaxCommissionChange.Cmp(node_manager.GenesisMaxCommissionChange) < 0 {
				globalConfig.MaxCommissionChange = config.MaxCommissionChange
			}
			if config.MinInitialStake.Sign() > 0 {
				globalConfig.MinInitialStake = config.MinInitialStake
			}
			if config.MinProposalStake.Sign() > 0 {
				globalConfig.MinProposalStake = config.MinProposalStake
			}
			err = node_manager.SetGlobalConfig(s, globalConfig)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, node_manager.SetGlobalConfig error: %v", err)
			}

			// remove from proposal list
			err = removeFromConfigProposalList(s, params.ID)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, removeFromConfigProposalList error: %v", err)
			}
		case UpdateCommunityInfo:
			info := new(community.CommunityInfo)
			err := rlp.DecodeBytes(proposal.Content, info)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, deserialize community info error: %v", err)
			}
			if info.CommunityAddress != common.EmptyAddress {
				communityInfo.CommunityAddress = info.CommunityAddress
			}
			if info.CommunityRate.Sign() > 0 {
				communityInfo.CommunityRate = info.CommunityRate
			}
			err = community.SetCommunityInfo(s, communityInfo)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, node_manager.SetCommunityInfo error: %v", err)
			}

			// remove from proposal list
			err = removeFromCommunityProposalList(s, params.ID)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, removeFromCommunityProposalList error: %v", err)
			}
		case UpdateSideChain:
			args := new(side_chain_manager.SideChain)
			err := rlp.DecodeBytes(proposal.Content, args)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, deserialize side chain params error: %v", err)
			}

			sideChainInfo, err := side_chain_manager.GetSideChainObject(s, args.ChainID)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, side_chain_manager.GetSideChainObject error: %v", err)
			}
			if sideChainInfo == nil {
				sideChainInfo = &side_chain_manager.SideChain{ChainID: args.ChainID}
			}
			if args.Router != 0 {
				sideChainInfo.Router = args.Router
			}
			if len(args.CCMCAddress) != 0 {
				sideChainInfo.CCMCAddress = args.CCMCAddress
			}
			if args.Name != "" {
				sideChainInfo.Name = args.Name
			}
			if len(args.ExtraInfo) != 0 {
				sideChainInfo.ExtraInfo = args.ExtraInfo
			}
			err = side_chain_manager.PutSideChain(s, sideChainInfo)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, side_chain_manager.PutSideChain error: %v", err)
			}

			// remove from proposal list
			err = removeFromSideChainProposalList(s, params.ID)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, removeFromSideChainProposalList error: %v", err)
			}
		case Normal:
			// remove from proposal list
			err = removeFromProposalList(s, params.ID)
			if err != nil {
				return nil, fmt.Errorf("VoteProposal, removeFromProposalList error: %v", err)
			}
		}

		err = s.AddNotify(ABI, []string{VOTE_PROPOSAL_EVENT}, proposal.ID.String())
		if err != nil {
			return nil, fmt.Errorf("VoteProposal, AddNotify error: %v", err)
		}
	}
	return utils.PackOutputs(ABI, MethodVoteProposal, true)
}

func GetProposal(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &GetProposalParam{}
	if err := utils.UnpackMethod(ABI, MethodGetProposal, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("VoteProposal, unpack params error: %v", err)
	}

	proposal, err := getProposal(s, params.ID)
	if err != nil {
		return nil, fmt.Errorf("GetProposal, getProposal error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(proposal)
	if err != nil {
		return nil, fmt.Errorf("GetProposal, serialize proposal error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetProposal, enc)
}

func GetProposalList(s *native.NativeContract) ([]byte, error) {
	proposalList, err := getProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("GetProposalList, getProposalList error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(proposalList)
	if err != nil {
		return nil, fmt.Errorf("GetProposalList, serialize proposal list error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetProposalList, enc)
}

func GetConfigProposalList(s *native.NativeContract) ([]byte, error) {
	configProposalList, err := getConfigProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("GetConfigProposalList, getConfigProposalList error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(configProposalList)
	if err != nil {
		return nil, fmt.Errorf("GetConfigProposalList, serialize config proposal list error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetConfigProposalList, enc)
}

func GetCommunityProposalList(s *native.NativeContract) ([]byte, error) {
	communityProposalList, err := getCommunityProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("GetCommunityProposalList, getCommunityProposalList error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(communityProposalList)
	if err != nil {
		return nil, fmt.Errorf("GetCommunityProposalList, serialize community proposal list error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetCommunityProposalList, enc)
}

func GetSideChainProposalList(s *native.NativeContract) ([]byte, error) {
	sideChainProposalList, err := getSideChainProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("GetSideChainProposalList, getSideChainProposalList error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(sideChainProposalList)
	if err != nil {
		return nil, fmt.Errorf("GetSideChainProposalList, serialize side chain proposal list error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetSideChainProposalList, enc)
}
