package cross_chain_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contract"
	"github.com/ethereum/go-ethereum/contract/utils"
	"github.com/ethereum/go-ethereum/modules/cfg"
)

func PutBlackChain(module *contract.ModuleContract, chainID uint64) error {
	module.GetCacheDB().Put(blackChainKey(chainID), utils.GetUint64Bytes(chainID))
	return nil
}

func RemoveBlackChain(module *contract.ModuleContract, chainID uint64) {
	module.GetCacheDB().Delete(blackChainKey(chainID))
}

func CheckIfChainBlacked(module *contract.ModuleContract, chainID uint64) (bool, error) {
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
	contractAddr := cfg.CrossChainManagerContractAddress
	return utils.ConcatKey(contractAddr, []byte(BLACKED_CHAIN), utils.GetUint64Bytes(chainID))
}
