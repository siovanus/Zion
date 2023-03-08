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

package boot

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/modules/cfg"
	"github.com/ethereum/go-ethereum/modules/cross_chain_manager"
	"github.com/ethereum/go-ethereum/modules/economic"
	"github.com/ethereum/go-ethereum/modules/info_sync"
	"github.com/ethereum/go-ethereum/modules/node_manager"
	"github.com/ethereum/go-ethereum/modules/proposal_manager"
	"github.com/ethereum/go-ethereum/modules/side_chain_manager"
)

func InitModuleContracts() {
	node_manager.InitNodeManager()
	economic.InitEconomic()
	info_sync.InitInfoSync()
	cross_chain_manager.InitCrossChainManager()
	side_chain_manager.InitSideChainManager()
	proposal_manager.InitProposalManager()

	log.Info("Initialize module contracts",
		"node manager", cfg.NodeManagerContractAddress.Hex(),
		"economic", cfg.EconomicContractAddress.Hex(),
		"header sync", cfg.InfoSyncContractAddress.Hex(),
		"cross chain manager", cfg.CrossChainManagerContractAddress.Hex(),
		"side chain manager", cfg.SideChainManagerContractAddress.Hex(),
		"proposal manager", cfg.ProposalManagerContractAddress.Hex(),
	)
}
