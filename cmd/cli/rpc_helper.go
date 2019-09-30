//   Copyright (C) 2018 XChain
//
//   This program is free software: you can redistribute it and/or modify
//   it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either cliVersion 3 of the License, or
//   (at your option) any later cliVersion.
//
//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//   GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
//   along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cli

import (
	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/core"
	"github.com/xchain/go-chain/global/types"
	"math/big"
)

func convertTransaction(tx *types.Transaction) *Transaction {
	var (
		gasLimit = uint64(0)
		gasPrice = uint64(0)
		value    = uint64(0)
	)
	if tx.GasLimit != nil {
		gasLimit = tx.GasLimit.Uint64()
	}
	if tx.GasPrice != nil {
		gasPrice = tx.GasPrice.Uint64()
	}
	if tx.Value != nil {
		value = tx.Value.Uint64()
	}
	trans := &Transaction{
		Hash:     tx.Hash,
		Source:   tx.Source,
		Target:   tx.Target,
		Type:     tx.Type,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Data:     tx.Data,
		Nonce:    tx.Nonce,
		Value:    common.AM2DDAM(value),
	}
	return trans
}

func convertExecutedTransaction(executed *types.ExecutedTransaction) *ExecutedTransaction {
	rec := &Receipt{
		Status:            int(executed.Receipt.Status),
		CumulativeGasUsed: executed.Receipt.CumulativeGasUsed,
		TxHash:            executed.Receipt.TxHash,
		Height:            executed.Receipt.Height,
		TxIndex:           executed.Receipt.TxIndex,
	}
	return &ExecutedTransaction{
		Receipt:     rec,
		Transaction: convertTransaction(executed.Transaction),
	}
}

func convertBlockHeader(b *types.Block) *Block {
	bh := b.Header
	diff := common.MaxUint64 / bh.BaseTarget
	block := &Block{
		Height:               bh.Height,
		Hash:                 bh.Hash,
		PreHash:              bh.PreHash,
		CurTime:              bh.CurTime.Local(),
		Proposer:             bh.Proposer,
		Nonce:                bh.Nonce,
		TxTree:               bh.TxTree,
		ReceiptTree:          bh.ReceiptTree,
		StateTree:            bh.StateTree,
		BaseTarget:           bh.BaseTarget,
		CumulativeDifficulty: bh.CumulativeDifficulty,
		Difficulty:           new(big.Int).SetUint64(diff.Uint64()),
		Capacity:             core.EstimateCapacity(bh.BaseTarget).String(),
		Auth:                 bh.Auth,
		Sign:                 bh.Sign,
	}
	return block
}
