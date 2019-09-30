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

package types

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/crypto"
	"github.com/xchain/go-chain/middleware/pb"
	time2 "github.com/xchain/go-chain/middleware/time"
	"math/big"
)

// UnMarshalTransactions deserialize from []byte to *Transaction
func UnMarshalTransactions(b []byte) ([]*Transaction, error) {
	ts := new(pb.TransactionSlice)
	error := proto.Unmarshal(b, ts)
	if error != nil {
		return nil, error
	}

	result := PbToTransactions(ts.Transactions)
	return result, nil
}

// UnMarshalBlock deserialize from []byte to *Block
func UnMarshalBlock(bytes []byte) (*Block, error) {
	b := new(pb.Block)
	error := proto.Unmarshal(bytes, b)
	if error != nil {
		return nil, error
	}
	block := PbToBlock(b)
	return block, nil
}

// UnMarshalBlockHeader deserialize from []byte to *BlockHeader
func UnMarshalBlockHeader(bytes []byte) (*BlockHeader, error) {
	b := new(pb.BlockHeader)
	error := proto.Unmarshal(bytes, b)
	if error != nil {
		return nil, error
	}
	header := PbToBlockHeader(b)
	return header, nil
}

// MarshalTransactions serialize []*Transaction
func MarshalTransactions(txs []*Transaction) ([]byte, error) {
	transactions := TransactionsToPb(txs)
	transactionSlice := pb.TransactionSlice{Transactions: transactions}
	return proto.Marshal(&transactionSlice)
}

// MarshalBlock serialize *Block
func MarshalBlock(b *Block) ([]byte, error) {
	block := BlockToPb(b)
	if block == nil {
		return nil, nil
	}
	return proto.Marshal(block)
}

// MarshalBlockHeader Serialize *BlockHeader
func MarshalBlockHeader(b *BlockHeader) ([]byte, error) {
	block := BlockHeaderToPb(b)
	if block == nil {
		return nil, nil
	}
	return proto.Marshal(block)
}

func ensureUint64(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
func ensureInt32(ptr *int32) int32 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
func ensureInt64(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func byteToHash(b []byte) common.Hash {
	if len(b) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(b)
}

func pbToTransaction(t *pb.Transaction) *Transaction {
	if t == nil {
		return &Transaction{}
	}

	var target *common.Address
	if t.Target != nil {
		t := common.BytesToAddress(t.Target)
		target = &t
	}
	value := NewBigInt(0).SetBytesWithSign(t.Value)
	gasLimit := NewBigInt(0).SetBytesWithSign(t.GasLimit)
	gasPrice := NewBigInt(0).SetBytesWithSign(t.GasPrice)

	transaction := &Transaction{
		Data:     t.Data,
		Value:    value,
		Nonce:    ensureUint64(t.Nonce),
		Target:   target,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Hash:     byteToHash(t.Hash),
		Type:     int8(ensureInt32(t.Type)),
		Sign:     crypto.BytesToSign(t.Sign),
	}
	return transaction
}

func PbToTransactions(txs []*pb.Transaction) []*Transaction {
	result := make([]*Transaction, 0)
	if txs == nil {
		return result
	}
	for _, t := range txs {
		transaction := pbToTransaction(t)
		result = append(result, transaction)
	}
	return result
}

func PbToBlockHeader(h *pb.BlockHeader) *BlockHeader {
	if h == nil {
		return nil
	}
	diff := new(big.Int)
	if h.CumulativeDifficulty != nil {
		diff = diff.SetBytes(h.CumulativeDifficulty)
	}

	header := BlockHeader{
		Hash:                 byteToHash(h.Hash),
		Height:               ensureUint64(h.Height),
		PreHash:              byteToHash(h.PreHash),
		CurTime:              time2.Int64ToTimeStamp(ensureInt64(h.CurTime)),
		Proposer:             common.BytesToAddress(h.Proposer),
		Nonce:                Nonce(ensureUint64(h.Nonce)),
		TxTree:               byteToHash(h.TxTree),
		ReceiptTree:          byteToHash(h.ReceiptTree),
		StateTree:            byteToHash(h.StateTree),
		BaseTarget:           Target(ensureUint64(h.BaseTarget)),
		CumulativeDifficulty: diff,
		Auth:                 h.Auth,
		Sign:                 crypto.BytesToSign(h.Sign),
	}
	return &header
}

func PbToBlock(b *pb.Block) *Block {
	if b == nil {
		return nil
	}
	h := PbToBlockHeader(b.Header)
	txs := PbToTransactions(b.Transactions)
	block := Block{Header: h, Transactions: txs}
	return &block
}

func transactionToPb(t *Transaction) *pb.Transaction {
	if t == nil {
		return nil
	}
	var (
		target []byte
	)
	if t.Target != nil {
		target = t.Target.Bytes()
	}

	tp := int32(t.Type)
	transaction := pb.Transaction{
		Data:     t.Data,
		Value:    t.Value.GetBytesWithSign(),
		Nonce:    &t.Nonce,
		Target:   target,
		GasLimit: t.GasLimit.GetBytesWithSign(),
		GasPrice: t.GasPrice.GetBytesWithSign(),
		Hash:     t.Hash.Bytes(),
		Type:     &tp,
		Sign:     t.Sign.Bytes(),
	}
	return &transaction
}

func TransactionsToPb(txs []*Transaction) []*pb.Transaction {
	if txs == nil {
		return nil
	}
	transactions := make([]*pb.Transaction, 0)
	for _, t := range txs {
		transaction := transactionToPb(t)
		transactions = append(transactions, transaction)
	}
	return transactions
}

func BlockHeaderToPb(h *BlockHeader) *pb.BlockHeader {
	ts := h.CurTime.Unix()
	nc := h.Nonce.Uint64()
	tg := h.BaseTarget.Uint64()
	header := pb.BlockHeader{
		Hash:                 h.Hash.Bytes(),
		Height:               &h.Height,
		PreHash:              h.PreHash.Bytes(),
		CurTime:              &ts,
		Proposer:             h.Proposer.Bytes(),
		Nonce:                &nc,
		TxTree:               h.TxTree.Bytes(),
		ReceiptTree:          h.ReceiptTree.Bytes(),
		StateTree:            h.StateTree.Bytes(),
		BaseTarget:           &tg,
		CumulativeDifficulty: h.CumulativeDifficulty.Bytes(),
		Auth:                 h.Auth,
		Sign:                 h.Sign.Bytes(),
	}
	return &header
}

func BlockToPb(b *Block) *pb.Block {
	if b == nil {
		return nil
	}
	header := BlockHeaderToPb(b.Header)
	transactions := TransactionsToPb(b.Transactions)
	block := pb.Block{Header: header, Transactions: transactions}
	return &block
}
