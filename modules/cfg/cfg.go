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

package cfg

import "github.com/ethereum/go-ethereum/common"

const (
	ModuleInfoSync         = "info_sync"
	ModuleCrossChain       = "cross_chain"
	ModuleNodeManager      = "node_manager"
	ModuleSideChainManager = "side_chain_manager"
	ModuleEconomic         = "economic"
	ModuleProposalManager  = "proposal_manager"
)

var (
	SystemTxSender = common.HexToAddress("0x0000000000000000000100000000000000000000")
)

var (
	NodeManagerContractAddress       = common.HexToAddress("0x0000000000000000000000000000000000001000")
	EconomicContractAddress          = common.HexToAddress("0x0000000000000000000000000000000000001001")
	InfoSyncContractAddress          = common.HexToAddress("0x0000000000000000000000000000000000001002")
	CrossChainManagerContractAddress = common.HexToAddress("0x0000000000000000000000000000000000001003")
	SideChainManagerContractAddress  = common.HexToAddress("0x0000000000000000000000000000000000001004")
	ProposalManagerContractAddress   = common.HexToAddress("0x0000000000000000000000000000000000001005")
)

var ModuleContractAddrMap = map[string]common.Address{
	ModuleNodeManager:      NodeManagerContractAddress,
	ModuleEconomic:         EconomicContractAddress,
	ModuleInfoSync:         InfoSyncContractAddress,
	ModuleCrossChain:       CrossChainManagerContractAddress,
	ModuleSideChainManager: SideChainManagerContractAddress,
	ModuleProposalManager:  ProposalManagerContractAddress,
}

func IsModuleContract(addr common.Address) bool {
	for _, v := range ModuleContractAddrMap {
		if v == addr {
			return true
		}
	}
	return false
}