package helper

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contract"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"math/big"
)

func ValidateOwner(n *contract.ModuleContract, address common.Address) error {
	if n.ContractRef().TxOrigin() != address {
		return fmt.Errorf("validateOwner, authentication failed!")
	}
	return nil
}

func ModuleTransfer(s *state.StateDB, from, to common.Address, amount *big.Int) error {
	if amount.Sign() == -1 {
		return fmt.Errorf("amount can not be negative")
	}
	if !core.CanTransfer(s, from, amount) {
		return fmt.Errorf("%s insufficient balance", from.Hex())
	}
	core.Transfer(s, from, to, amount)
	return nil
}