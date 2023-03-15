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

package info_sync

import (
	"fmt"

	"github.com/ethereum/go-ethereum/contract"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/modules/cfg"
	"github.com/ethereum/go-ethereum/modules/node_manager"
	"github.com/ethereum/go-ethereum/modules/side_chain_manager"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	this = cfg.ModuleContractAddrMap[cfg.ModuleInfoSync]
)

func InitInfoSync() {
	ABI = GetABI()
	contract.Contracts.RegisterContract(this, RegisterInfoSyncContract)
}

func RegisterInfoSyncContract(s *contract.ModuleContract) {
	s.Register(MethodContractName, Name)
	s.Register(MethodSyncRootInfo, SyncRootInfo)
	s.Register(MethodReplenish, Replenish)
	s.Register(MethodGetInfoHeight, GetInfoHeight)
	s.Register(MethodGetInfo, GetInfo)
}

func Name(s *contract.ModuleContract) ([]byte, error) {
	return contract.PackOutputs(ABI, MethodContractName, cfg.ModuleInfoSync)
}

func SyncRootInfo(s *contract.ModuleContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &SyncRootInfoParam{}
	if err := contract.UnpackMethod(ABI, MethodSyncRootInfo, params, ctx.Payload); err != nil {
		return nil, err
	}
	chainID := params.ChainID

	//check if chainid exist
	sideChain, err := side_chain_manager.GetSideChainObject(s, chainID)
	if err != nil {
		return nil, fmt.Errorf("SyncRootInfo, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("SyncRootInfo, side chain is not registered")
	}

	//verify signature
	digest, err := params.Digest()
	if err != nil {
		return nil, fmt.Errorf("SyncRootInfo, digest input param error: %v", err)
	}
	pub, err := crypto.SigToPub(digest, params.Signature)
	if err != nil {
		return nil, fmt.Errorf("SyncRootInfo, crypto.SigToPub error: %v", err)
	}
	addr := crypto.PubkeyToAddress(*pub)

	//sync root infos
	for _, v := range params.RootInfos {
		var rootInfo *RootInfo
		err := rlp.DecodeBytes(v, &rootInfo)
		if err != nil {
			return nil, fmt.Errorf("SyncRootInfo, decode root info error")
		}
		//use chain id, info key and value as unique id
		unique := &RootInfoUnique{
			ChainID: params.ChainID,
			Height:  rootInfo.Height,
			Info:    rootInfo.Info,
		}
		blob, err := rlp.EncodeToBytes(unique)
		if err != nil {
			return nil, err
		}

		ok, err := node_manager.CheckConsensusSigns(s, MethodSyncRootInfo, blob, addr, node_manager.Voter)
		if err != nil {
			return nil, fmt.Errorf("SyncRootInfo, CheckVoterSigns error: %v", err)
		}
		if ok {
			err := PutRootInfo(s, chainID, rootInfo.Height, rootInfo.Info)
			if err != nil {
				return nil, fmt.Errorf("SyncRootInfo, PutCrossChainInfo error: %v", err)
			}
		}
	}

	return contract.PackOutputs(ABI, MethodSyncRootInfo, true)
}

func Replenish(s *contract.ModuleContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ReplenishParam{}
	if err := contract.UnpackMethod(ABI, MethodReplenish, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("Replenish, unpack params error: %s", err)
	}

	err := NotifyReplenish(s, params.Heights, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("Replenish, NotifyReplenish error: %s", err)
	}
	return contract.PackOutputs(ABI, MethodReplenish, true)
}

func GetInfoHeight(s *contract.ModuleContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &GetInfoHeightParam{}
	if err := contract.UnpackMethod(ABI, MethodGetInfoHeight, params, ctx.Payload); err != nil {
		return nil, err
	}

	height, err := GetCurrentHeight(s, params.ChainID)
	if err != nil {
		return nil, err
	}
	return contract.PackOutputs(ABI, MethodGetInfoHeight, height)
}

func GetInfo(s *contract.ModuleContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &GetInfoParam{}
	if err := contract.UnpackMethod(ABI, MethodGetInfo, params, ctx.Payload); err != nil {
		return nil, err
	}
	info, err := GetRootInfo(s, params.ChainID, params.Height)
	if err != nil {
		return nil, err
	}
	return contract.PackOutputs(ABI, MethodGetInfo, info)
}
