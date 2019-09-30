//   Copyright (C) 2018 XChain
//
//   This program is free software: you can redistribute it and/or modify
//   it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either version 3 of the License, or
//   (at your option) any later version.
//
//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//   GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
//   along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package types define the key data structures for the core
package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/xchain/go-chain/crypto"

	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/middleware/time"
)

const (
	txFixSize = 200 // Fixed size for each transaction
)

// Supported transaction types
const (
	TransactionTypeTransfer      = 0
	TransactionTypeBindUMID      = 1
	TransactionTypeTransformUMID = 2
	TransactionTypeUnbindUMID    = 3
	TransactionTypeStakeAdd      = 4
	TransactionTypeStakeReduce   = 5
)

// Transaction denotes one transaction infos
type Transaction struct {
	Data   []byte          `msgpack:"dt,omitempty"` // Data of the transaction, cost gas
	Value  *BigInt         `msgpack:"v"`            // The value the sender suppose to transfer
	Nonce  uint64          `msgpack:"nc"`           // The nonce indicates the transaction sequence related to sender
	Target *common.Address `msgpack:"tg,omitempty"` // The receiver address
	Type   int8            `msgpack:"tp"`           // Transaction type

	GasLimit *BigInt     `msgpack:"gl"`
	GasPrice *BigInt     `msgpack:"gp"`
	Hash     common.Hash `msgpack:"h"`

	Sign   *crypto.Sign    `msgpack:"si"`  // The Sign of the sender
	Source *common.Address `msgpack:"src"` // Sender address, recovered from sign
}

// GenHash generate unique hash of the transaction. source,sign is out of the hash calculation range
func (tx *Transaction) GenHash() common.Hash {
	if nil == tx {
		return common.Hash{}
	}
	buffer := bytes.Buffer{}
	if tx.Data != nil {
		buffer.Write(tx.Data)
	}
	buffer.Write(common.Uint64ToByte(tx.Nonce))
	if tx.Target != nil {
		buffer.Write(tx.Target.Bytes())
	}
	buffer.WriteByte(byte(tx.Type))
	if tx.Value != nil {
		buffer.Write(tx.Value.GetBytesWithSign())
	}
	buffer.Write(tx.GasLimit.GetBytesWithSign())
	buffer.Write(tx.GasPrice.GetBytesWithSign())

	return common.BytesToHash(common.Sha256(buffer.Bytes()))
}

// RecoverSource recover source from the sign field.
// It returns directly if source is not nil or it is a reward transaction.
func (tx *Transaction) RecoverSource() error {
	if tx.Source != nil {
		return nil
	}
	if tx.Sign == nil {
		return fmt.Errorf("sign is nil")
	}
	pk, err := tx.Sign.RecoverPubkey(tx.Hash.Bytes())
	if err == nil {
		src := pk.GetAddress()
		tx.Source = &src
	}
	return err
}

func (tx *Transaction) Size() int {
	return txFixSize + len(tx.Data)
}

func (tx Transaction) GetData() []byte { return tx.Data }

func (tx Transaction) GetGasLimit() uint64 {
	return tx.GasLimit.Uint64()
}
func (tx Transaction) GetValue() uint64 {
	return tx.Value.Uint64()
}

func (tx Transaction) GetSource() *common.Address { return tx.Source }
func (tx Transaction) GetTarget() *common.Address { return tx.Target }
func (tx Transaction) GetHash() common.Hash       { return tx.Hash }

type Target uint64

func (d Target) Bytes() []byte {
	return common.Uint64ToByte(uint64(d))
}

func (d Target) Uint64() uint64 {
	return uint64(d)
}

type Nonce uint64

func (d Nonce) Bytes() []byte {
	return common.Uint64ToByte(uint64(d))
}
func (d Nonce) Size() int {
	return 8
}
func (d Nonce) Uint64() uint64 {
	return uint64(d)
}

// BlockHeader is block header structure
type BlockHeader struct {
	Hash                 common.Hash    // The hash of this block
	Height               uint64         // The height of this block
	PreHash              common.Hash    // The hash of previous block
	CurTime              time.TimeStamp // Current block time
	Proposer             common.Address // Proposer Address
	Nonce                Nonce          // Nonce of the scoop
	TxTree               common.Hash    // Transaction Merkel root hash
	ReceiptTree          common.Hash    // Receipt Merkel root hash
	StateTree            common.Hash    // State db Merkel root hash
	BaseTarget           Target         // Target difficulty
	CumulativeDifficulty *big.Int
	Auth                 AuthCode
	Sign                 *crypto.Sign
}

func (bh *BlockHeader) Weight() *big.Int {
	return bh.CumulativeDifficulty
}

func (bh *BlockHeader) MoreWeightThan(chain ChainWeight) bool {
	return bh.Weight().Cmp(chain.Weight()) > 0
}

// GenHash calculates the hash of the block
func (bh *BlockHeader) GenHash() common.Hash {
	buf := bytes.NewBuffer([]byte{})

	buf.Write(common.UInt64ToByte(bh.Height))

	buf.Write(bh.PreHash.Bytes())

	buf.Write(bh.CurTime.Bytes())
	buf.Write(bh.Proposer.Bytes())
	buf.Write(bh.Nonce.Bytes())

	buf.Write(bh.TxTree.Bytes())
	buf.Write(bh.ReceiptTree.Bytes())
	buf.Write(bh.StateTree.Bytes())
	buf.Write(bh.BaseTarget.Bytes())
	buf.Write(bh.CumulativeDifficulty.Bytes())
	buf.Write(bh.Auth.Bytes())
	return common.BytesToHash(common.Sha256(buf.Bytes()))
}

func (bh *BlockHeader) HasTransactions() bool {
	return bh.TxTree != common.EmptyHash
}

// Block is the block data structure consists of the header and transactions as body
type Block struct {
	Header       *BlockHeader
	Transactions []*Transaction
}

func (b *Block) GetTransactionHashs() []common.Hash {
	if b.Transactions == nil {
		return []common.Hash{}
	}
	hashs := make([]common.Hash, 0)
	for _, tx := range b.Transactions {
		hashs = append(hashs, tx.Hash)
	}
	return hashs
}
