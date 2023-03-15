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

package state

import (
	"fmt"
	"math"
)

type GasMeter struct {
	gasLeft *uint64
}

func NewGasMeter(gasLeft *uint64) *GasMeter {
	return &GasMeter{gasLeft}
}

// only used when you need not to pay gas fee
func NewInfiniteGasMeter() *GasMeter {
	var infinite uint64 = math.MaxUint64
	return &GasMeter{&infinite}
}

func (gm *GasMeter) ConsumeGas(gasUsage uint64) error {
	if *gm.gasLeft < gasUsage {
		*gm.gasLeft = 0
		return fmt.Errorf("gasLeft not enough, need %d, got %d", gasUsage, *gm.gasLeft)
	}
	*gm.gasLeft = *gm.gasLeft - gasUsage
	return nil
}

type GasConfig struct {
	DeleteCost       uint64
	ReadCostFlat     uint64
	ReadCostPerByte  uint64
	WriteCostFlat    uint64
	WriteCostPerByte uint64
}

// DefaultGasConfig returns a default gas config for KVStores.
func DefaultGasConfig() *GasConfig {
	return &GasConfig{
		DeleteCost:       1000,
		ReadCostFlat:     1000,
		ReadCostPerByte:  3,
		WriteCostFlat:    2000,
		WriteCostPerByte: 30,
	}
}
