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
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrKeyLen       = errors.New("cacheDB should only be used for native contract storage")
	ErrBigLen       = errors.New("big int length out of range")
	ErrNotExist     = errors.New("key not exist")
	ErrInvalidBytes = errors.New("invalid bytes")
)

type Store struct {
	db        *StateDB
	gasMeter  *GasMeter
	gasConfig *GasConfig
}

func NewStore(db *StateDB, gasMeter *GasMeter, gasConfig *GasConfig) *Store {
	return &Store{db, gasMeter, gasConfig}
}

func (s *Store) Db() *StateDB {
	return s.db
}

func (s *Store) GasMeter() *GasMeter {
	return s.gasMeter
}

func (s *Store) GasConfig() *GasConfig {
	return s.gasConfig
}

// support storage for type of `Address`
func (s *Store) SetAddress(key []byte, value common.Address) error {
	hash := common.BytesToHash(value.Bytes())
	_, _, err := s.customSet(key, hash)
	return err
}

func (s *Store) GetAddress(key []byte) (common.Address, error) {
	_, _, hash, err := s.customGet(key)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(hash.Bytes()), nil
}

func (s *Store) DelAddress(key []byte) error {
	_, _, err := s.customDel(key)
	return err
}

// support storage for type of `Hash`
func (s *Store) SetHash(key []byte, value common.Hash) error {
	_, _, err := s.customSet(key, value)
	return err
}

func (s *Store) GetHash(key []byte) (common.Hash, error) {
	_, _, value, err := s.customGet(key)
	return value, err
}

func (s *Store) DelHash(key []byte) error {
	_, _, err := s.customDel(key)
	return err
}

// support storage for type of `big`
func (s *Store) SetBigInt(key []byte, value *big.Int) error {
	if len(value.Bytes()) > common.HashLength {
		return ErrBigLen
	}
	hash := common.BytesToHash(value.Bytes())
	_, _, err := s.customSet(key, hash)
	return err
}

func (s *Store) GetBigInt(key []byte) (*big.Int, error) {
	_, _, raw, err := s.customGet(key)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(raw[:]), nil
}

func (s *Store) DelBigInt(key []byte) error {
	_, _, err := s.customDel(key)
	return err
}

func (s *Store) SetBytes(key []byte, value []byte) error {
	s.Put(key, value)
	return nil
}

func (s *Store) GetBytes(key []byte) ([]byte, error) {
	return s.Get(key)
}

// custom functions
func (s *Store) customSet(key []byte, value common.Hash) (addr common.Address, slot common.Hash, err error) {
	addr, slot, err = parseKey(key)
	if err != nil {
		return
	}

	err = s.gasMeter.ConsumeGas(s.gasConfig.WriteCostFlat)
	if err != nil {
		return
	}
	s.db.SetState(addr, slot, value)

	err = s.gasMeter.ConsumeGas(s.gasConfig.WriteCostPerByte * uint64(len(value)))
	if err != nil {
		return
	}
	return
}

func (s *Store) customGet(key []byte) (addr common.Address, slot, value common.Hash, err error) {
	addr, slot, err = parseKey(key)
	if err != nil {
		return
	}

	err = s.gasMeter.ConsumeGas(s.gasConfig.ReadCostFlat)
	if err != nil {
		return
	}
	value = s.db.GetState(addr, slot)

	err = s.gasMeter.ConsumeGas(s.gasConfig.ReadCostPerByte * uint64(len(value)))
	if err != nil {
		return
	}
	return
}

func (s *Store) customDel(key []byte) (addr common.Address, slot common.Hash, err error) {
	addr, slot, err = parseKey(key)
	if err != nil {
		return
	}

	err = s.gasMeter.ConsumeGas(s.gasConfig.DeleteCost)
	if err != nil {
		return
	}
	s.db.SetState(addr, slot, common.Hash{})
	return
}

func parseKey(key []byte) (addr common.Address, slot common.Hash, err error) {
	if len(key) <= common.AddressLength {
		return common.Address{}, common.Hash{}, ErrKeyLen
	}
	addr = common.BytesToAddress(key[:common.AddressLength])
	slot = Key2Slot(key[common.AddressLength:])
	return
}

func (s *Store) Put(key []byte, value []byte) error {
	if len(key) <= common.AddressLength {
		panic("CacheDB should only be used for native contract storage")
	}

	err := s.Delete(key)
	if err != nil {
		return err
	}

	err = s.gasMeter.ConsumeGas(s.gasConfig.WriteCostFlat)
	if err != nil {
		return err
	}
	so := s.db.GetOrNewStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if so != nil {
		slot := Key2Slot(key[common.AddressLength:])
		if len(value) <= common.HashLength-1 {
			s.putValue(so, slot, value, false)
			value = nil
		} else {
			s.putValue(so, slot, value[:common.HashLength-1], true)
			value = value[common.HashLength-1:]
		}

		for len(value) > 0 {
			slot = s.nextSlot(slot)
			if len(value) <= common.HashLength-1 {
				s.putValue(so, slot, value, false)
				break
			} else {
				s.putValue(so, slot, value[:common.HashLength-1], true)
				value = value[common.HashLength-1:]
			}
		}

		err := s.gasMeter.ConsumeGas(s.gasConfig.WriteCostPerByte * uint64(len(value)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) putValue(so *stateObject, slot common.Hash, value []byte, more bool) {
	if len(value) > common.HashLength-1 {
		panic("value should not exceed 31")
	}

	if more && len(value) != common.HashLength-1 {
		panic("value length should equal 31 when more is true")
	}

	if more {
		value = append([]byte{1}, value...)
	} else {
		padding := make([]byte, common.HashLength-len(value))
		padding[0] = byte(len(value) << 1)
		value = append(padding, value...)
	}

	hashValue := common.BytesToHash(value)
	so.SetState(s.db.db, slot, hashValue)
}

func Key2Slot(key []byte) common.Hash {
	key = crypto.Keccak256(key)
	return common.BytesToHash(key)
}

func (s *Store) nextSlot(slot common.Hash) common.Hash {
	slotBytes := slot.Bytes()
	for offset := common.HashLength - 1; offset >= 0; offset-- {
		slotBytes[offset] = slotBytes[offset] + 1
		if slotBytes[offset] != 0 {
			break
		}
	}

	return Key2Slot(slotBytes)
}

func (s *Store) Get(key []byte) ([]byte, error) {
	if len(key) <= common.AddressLength {
		panic("CacheDB should only be used for native contract storage")
	}

	err := s.gasMeter.ConsumeGas(s.gasConfig.ReadCostFlat)
	if err != nil {
		return nil, err
	}
	so := s.db.getStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if so != nil {
		var result []byte
		slot := Key2Slot(key[common.AddressLength:])
		value := so.GetState(s.db.db, slot)
		meta := value[:][0]
		more := meta&1 == 1
		if more {
			result = append(result, value[1:]...)
		} else {
			if value == (common.Hash{}) {
				return nil, nil
			}
			result = append(result, value[common.HashLength-meta>>1:]...)
		}

		for more {
			slot = s.nextSlot(slot)
			value = so.GetState(s.db.db, slot)
			meta = value[:][0]
			more = meta&1 == 1
			if more {
				result = append(result, value[1:]...)
			} else {
				result = append(result, value[common.HashLength-meta>>1:]...)
			}
		}

		err = s.gasMeter.ConsumeGas(s.gasConfig.ReadCostPerByte * uint64(len(result)))
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, nil
}

func (s *Store) Delete(key []byte) error {
	if len(key) <= common.AddressLength {
		panic("CacheDB store should only be used for native contract storage")
	}

	err := s.gasMeter.ConsumeGas(s.gasConfig.DeleteCost)
	if err != nil {
		return err
	}
	so := s.db.GetOrNewStateObject(common.BytesToAddress(key[:common.AddressLength]))
	if so != nil {
		slot := Key2Slot(key[common.AddressLength:])
		value := so.GetState(s.db.db, slot)
		so.SetState(s.db.db, slot, common.Hash{})
		more := value[:][0]&1 == 1
		for more {
			slot = s.nextSlot(slot)
			value = so.GetState(s.db.db, slot)
			so.SetState(s.db.db, slot, common.Hash{})
			more = value[:][0]&1 == 1
		}
	}
	return nil
}
