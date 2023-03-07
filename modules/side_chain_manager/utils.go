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

package side_chain_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/modules"
	utils2 "github.com/ethereum/go-ethereum/modules/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/rlp"
)

func GetSideChainApply(module *modules.ModuleContract, chanid uint64) (*SideChain, error) {
	contract := utils2.SideChainManagerContractAddress
	chainidByte := utils2.GetUint64Bytes(chanid)

	sideChainStore, err := module.GetCacheDB().Get(utils2.ConcatKey(contract, []byte(SIDE_CHAIN_APPLY),
		chainidByte))
	if err != nil {
		return nil, fmt.Errorf("getRegisterSideChain,get registerSideChainRequestStore error: %v", err)
	}
	sideChain := new(SideChain)
	if sideChainStore != nil {
		if err := rlp.DecodeBytes(sideChainStore, sideChain); err != nil {
			return nil, fmt.Errorf("getRegisterSideChain, deserialize sideChain error: %v", err)
		}
		return sideChain, nil
	} else {
		return nil, nil
	}
}

func putSideChainApply(module *modules.ModuleContract, sideChain *SideChain) error {
	contract := utils2.SideChainManagerContractAddress
	chainidByte := utils2.GetUint64Bytes(sideChain.ChainID)

	blob, err := rlp.EncodeToBytes(sideChain)
	if err != nil {
		return fmt.Errorf("putRegisterSideChain, sideChain.Serialization error: %v", err)
	}

	module.GetCacheDB().Put(utils2.ConcatKey(contract, []byte(SIDE_CHAIN_APPLY), chainidByte), blob)
	return nil
}

func GetSideChainObject(module *modules.ModuleContract, chainID uint64) (*SideChain, error) {
	contract := utils2.SideChainManagerContractAddress
	chainIDByte := utils2.GetUint64Bytes(chainID)

	sideChainStore, err := module.GetCacheDB().Get(utils2.ConcatKey(contract, []byte(SIDE_CHAIN),
		chainIDByte))
	if err != nil {
		return nil, fmt.Errorf("getSideChain,get registerSideChainRequestStore error: %v", err)
	}
	sideChain := new(SideChain)
	if sideChainStore != nil {
		if err := rlp.DecodeBytes(sideChainStore, sideChain); err != nil {
			return nil, fmt.Errorf("getSideChain, deserialize sideChain error: %v", err)
		}
		return sideChain, nil
	}
	return nil, nil

}

func PutSideChain(module *modules.ModuleContract, sideChain *SideChain) error {
	contract := utils2.SideChainManagerContractAddress
	chainidByte := utils2.GetUint64Bytes(sideChain.ChainID)

	blob, err := rlp.EncodeToBytes(sideChain)
	if err != nil {
		return fmt.Errorf("putSideChain, sideChain.Serialization error: %v", err)
	}

	module.GetCacheDB().Put(utils2.ConcatKey(contract, []byte(SIDE_CHAIN), chainidByte), blob)
	return nil
}

func getUpdateSideChain(module *modules.ModuleContract, chanid uint64) (*SideChain, error) {
	contract := utils2.SideChainManagerContractAddress
	chainidByte := utils2.GetUint64Bytes(chanid)

	sideChainStore, err := module.GetCacheDB().Get(utils2.ConcatKey(contract, []byte(UPDATE_SIDE_CHAIN_REQUEST),
		chainidByte))
	if err != nil {
		return nil, fmt.Errorf("getUpdateSideChain,get registerSideChainRequestStore error: %v", err)
	}
	sideChain := new(SideChain)
	if sideChainStore != nil {
		if err := rlp.DecodeBytes(sideChainStore, sideChain); err != nil {
			return nil, fmt.Errorf("getUpdateSideChain, deserialize sideChain error: %v", err)
		}
		return sideChain, nil
	} else {
		return nil, nil
	}
}

func putUpdateSideChain(module *modules.ModuleContract, sideChain *SideChain) error {
	contract := utils2.SideChainManagerContractAddress
	chainidByte := utils2.GetUint64Bytes(sideChain.ChainID)

	blob, err := rlp.EncodeToBytes(sideChain)
	if err != nil {
		return fmt.Errorf("putUpdateSideChain, sideChain.Serialization error: %v", err)
	}

	module.GetCacheDB().Put(utils2.ConcatKey(contract, []byte(UPDATE_SIDE_CHAIN_REQUEST), chainidByte), blob)
	return nil
}

func getQuitSideChain(module *modules.ModuleContract, chainid uint64) error {
	contract := utils2.SideChainManagerContractAddress
	chainidByte := utils2.GetUint64Bytes(chainid)

	chainIDStore, err := module.GetCacheDB().Get(utils2.ConcatKey(contract, []byte(QUIT_SIDE_CHAIN_REQUEST),
		chainidByte))
	if err != nil {
		return fmt.Errorf("getQuitSideChain, get registerSideChainRequestStore error: %v", err)
	}
	if chainIDStore != nil {
		return nil
	}
	return fmt.Errorf("getQuitSideChain, no record")
}

func putQuitSideChain(module *modules.ModuleContract, chainid uint64) error {
	contract := utils2.SideChainManagerContractAddress
	chainidByte := utils2.GetUint64Bytes(chainid)

	module.GetCacheDB().Put(utils2.ConcatKey(contract, []byte(QUIT_SIDE_CHAIN_REQUEST), chainidByte), chainidByte)
	return nil
}

func PutFee(module *modules.ModuleContract, chainId uint64, fee *Fee) error {
	contract := utils2.SideChainManagerContractAddress
	chainIdBytes := utils2.GetUint64Bytes(chainId)
	blob, err := rlp.EncodeToBytes(fee)
	if err != nil {
		return fmt.Errorf("PutFee, rlp.EncodeToBytes fee error: %v", err)
	}
	module.GetCacheDB().Put(utils2.ConcatKey(contract, []byte(FEE), chainIdBytes), blob)
	return nil
}

func GetFeeObj(module *modules.ModuleContract, chainID uint64) (*Fee, error) {
	chainIDBytes := utils2.GetUint64Bytes(chainID)
	key := utils2.ConcatKey(utils2.SideChainManagerContractAddress, []byte(FEE), chainIDBytes)
	store, err := module.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("GetFee, get fee info store error: %v", err)
	}
	fee := &Fee{
		Fee: new(big.Int),
	}
	if store != nil {
		if err := rlp.DecodeBytes(store, fee); err != nil {
			return nil, fmt.Errorf("GetFee, deserialize fee error: %v", err)
		}
	}
	return fee, nil
}

func PutFeeInfo(module *modules.ModuleContract, chainId, view uint64, feeInfo *FeeInfo) error {
	chainIdBytes := utils2.GetUint64Bytes(chainId)
	viewBytes := utils2.GetUint64Bytes(view)
	key := utils2.ConcatKey(utils2.SideChainManagerContractAddress, []byte(FEE_INFO), chainIdBytes, viewBytes)
	blob, err := rlp.EncodeToBytes(feeInfo)
	if err != nil {
		return fmt.Errorf("PutFeeInfo, rlp.EncodeToBytes fee info error: %v", err)
	}
	module.GetCacheDB().Put(key, blob)
	return nil
}

func GetFeeInfo(module *modules.ModuleContract, chainID, view uint64) (*FeeInfo, error) {
	chainIDBytes := utils2.GetUint64Bytes(chainID)
	viewBytes := utils2.GetUint64Bytes(view)
	key := utils2.ConcatKey(utils2.SideChainManagerContractAddress, []byte(FEE_INFO), chainIDBytes, viewBytes)
	store, err := module.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("GetFeeInfo, get fee info store error: %v", err)
	}
	feeInfo := &FeeInfo{
		FeeInfo: make(map[common.Address]*big.Int),
	}
	if store != nil {
		if err := rlp.DecodeBytes(store, feeInfo); err != nil {
			return nil, fmt.Errorf("GetFeeInfo, deserialize fee info error: %v", err)
		}
	}
	return feeInfo, nil
}

func GetRippleExtraInfo(module *modules.ModuleContract, chainId uint64) (*RippleExtraInfo, error) {
	sideChainInfo, err := GetSideChainObject(module, chainId)
	if err != nil {
		return nil, fmt.Errorf("GetRippleExtraInfo, GetSideChainObject error: %v", err)
	}
	if sideChainInfo == nil {
		return nil, fmt.Errorf("GetRippleExtraInfo, side chain info is nil")
	}
	rippleExtraInfo := &RippleExtraInfo{
		Pks: make([][]byte, 0),
	}
	if err := rlp.DecodeBytes(sideChainInfo.ExtraInfo, rippleExtraInfo); err != nil {
		return nil, fmt.Errorf("GetRippleExtraInfo, deserialize info error: %v", err)
	}
	return rippleExtraInfo, nil
}

func PutRippleExtraInfo(module *modules.ModuleContract, chainId uint64, rippleExtraInfo *RippleExtraInfo) error {
	blob, err := rlp.EncodeToBytes(rippleExtraInfo)
	if err != nil {
		return fmt.Errorf("PutRippleExtraInfo, rlp.EncodeToBytes info error: %v", err)
	}
	sideChainInfo, err := GetSideChainObject(module, chainId)
	if err != nil {
		return fmt.Errorf("PutRippleExtraInfo, GetSideChainObject error: %v", err)
	}
	sideChainInfo.ExtraInfo = blob
	err = PutSideChain(module, sideChainInfo)
	if err != nil {
		return fmt.Errorf("PutRippleExtraInfo, PutSideChain error: %v", err)
	}
	return nil
}

func PutAssetBind(module *modules.ModuleContract, chainId uint64, assetBind *AssetBind) error {
	chainIDBytes := utils2.GetUint64Bytes(chainId)
	key := utils2.ConcatKey(utils2.SideChainManagerContractAddress, []byte(ASSET_BIND), chainIDBytes)
	blob, err := rlp.EncodeToBytes(assetBind)
	if err != nil {
		return fmt.Errorf("PutAssetBind, rlp.EncodeToBytes asset bind error: %v", err)
	}
	module.GetCacheDB().Put(key, blob)
	return nil
}

func GetAssetBind(module *modules.ModuleContract, chainId uint64) (*AssetBind, error) {
	chainIDBytes := utils2.GetUint64Bytes(chainId)
	key := utils2.ConcatKey(utils2.SideChainManagerContractAddress, []byte(ASSET_BIND), chainIDBytes)
	store, err := module.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("GetAssetBind, get asset map store error: %v", err)
	}
	assetBind := &AssetBind{
		AssetMap:     make(map[uint64][]byte),
		LockProxyMap: make(map[uint64][]byte),
	}
	if store != nil {
		if err := rlp.DecodeBytes(store, assetBind); err != nil {
			return nil, fmt.Errorf("GetAssetBind, deserialize info error: %v", err)
		}
	}
	return assetBind, nil
}
