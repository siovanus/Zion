package info_sync

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contract"
	"github.com/ethereum/go-ethereum/modules/cfg"
)

const (
	//key prefix
	ROOT_INFO            = "rootInfo"
	CURRENT_HEIGHT       = "currentHeight"
	SYNC_ROOT_INFO_EVENT = "SyncRootInfoEvent"
	REPLENISH_EVENT      = "ReplenishEvent"
)

func PutRootInfo(module *contract.ModuleContract, chainID uint64, height uint32, info []byte) error {
	contractAddr := cfg.InfoSyncContractAddress
	chainIDBytes := contract.GetUint64Bytes(chainID)
	heightBytes := contract.GetUint32Bytes(height)

	module.GetCacheDB().Put(contract.ConcatKey(contractAddr, []byte(ROOT_INFO), chainIDBytes, heightBytes),
		info)
	currentHeight, err := GetCurrentHeight(module, chainID)
	if err != nil {
		return fmt.Errorf("PutRootInfo, GetCurrentHeight error: %v", err)
	}
	if currentHeight < height {
		module.GetCacheDB().Put(contract.ConcatKey(contractAddr, []byte(CURRENT_HEIGHT), chainIDBytes), heightBytes)
	}
	err = NotifyPutRootInfo(module, chainID, height)
	if err != nil {
		return fmt.Errorf("PutRootInfo, NotifyPutRootInfo error: %v", err)
	}
	return nil
}

func GetRootInfo(module *contract.ModuleContract, chainID uint64, height uint32) ([]byte, error) {
	contractAddr := cfg.InfoSyncContractAddress
	chainIDBytes := contract.GetUint64Bytes(chainID)
	heightBytes := contract.GetUint32Bytes(height)

	r, err := module.GetCacheDB().Get(contract.ConcatKey(contractAddr, []byte(ROOT_INFO), chainIDBytes, heightBytes))
	if err != nil {
		return nil, fmt.Errorf("GetRootInfo, module.GetCacheDB().Get error: %v", err)
	}
	return r, nil
}

func GetCurrentHeight(module *contract.ModuleContract, chainID uint64) (uint32, error) {
	contractAddr := cfg.InfoSyncContractAddress
	chainIDBytes := contract.GetUint64Bytes(chainID)

	r, err := module.GetCacheDB().Get(contract.ConcatKey(contractAddr, []byte(CURRENT_HEIGHT), chainIDBytes))
	if err != nil {
		return 0, fmt.Errorf("GetCurrentHeight, module.GetCacheDB().Get error: %v", err)
	}
	return contract.GetBytesUint32(r), nil
}

func NotifyPutRootInfo(module *contract.ModuleContract, chainID uint64, height uint32) error {
	err := module.AddNotify(ABI, []string{SYNC_ROOT_INFO_EVENT}, chainID, height, module.ContractRef().BlockHeight())
	if err != nil {
		return fmt.Errorf("NotifyPutRootInfo failed: %v", err)
	}
	return nil
}

func NotifyReplenish(module *contract.ModuleContract, heights []uint32, chainId uint64) error {
	err := module.AddNotify(ABI, []string{REPLENISH_EVENT}, heights, chainId)
	if err != nil {
		return fmt.Errorf("NotifyReplenish failed: %v", err)
	}
	return nil
}
