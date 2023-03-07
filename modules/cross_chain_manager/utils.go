package cross_chain_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/modules"
	utils2 "github.com/ethereum/go-ethereum/modules/utils"
)

func PutBlackChain(module *modules.ModuleContract, chainID uint64) error {
	module.GetCacheDB().Put(blackChainKey(chainID), utils2.GetUint64Bytes(chainID))
	return nil
}

func RemoveBlackChain(module *modules.ModuleContract, chainID uint64) {
	module.GetCacheDB().Delete(blackChainKey(chainID))
}

func CheckIfChainBlacked(module *modules.ModuleContract, chainID uint64) (bool, error) {
	chainIDStore, err := module.GetCacheDB().Get(blackChainKey(chainID))
	if err != nil {
		return true, fmt.Errorf("CheckBlackChain, get black chainIDStore error: %v", err)
	}
	if chainIDStore == nil {
		return false, nil
	}
	return true, nil
}

func blackChainKey(chainID uint64) []byte {
	contract := utils2.CrossChainManagerContractAddress
	return utils2.ConcatKey(contract, []byte(BLACKED_CHAIN), utils2.GetUint64Bytes(chainID))
}
