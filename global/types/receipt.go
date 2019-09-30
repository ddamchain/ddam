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
	"unsafe"

	"github.com/xchain/go-chain/common"
)

//go:generate gencodec -type Receipt -field-override receiptMarshaling -out gen_receipt_json.go

type ReceiptStatus int

const (
	RSSuccess ReceiptStatus = iota
	RSFail
	RSBalanceNotEnough
	RSParseFail
)

type Receipt struct {
	PostState         []byte        `json:"-"`
	Status            ReceiptStatus `json:"status"`
	CumulativeGasUsed uint64        `json:"cumulativeGasUsed"`

	TxHash  common.Hash `json:"transactionHash" gencodec:"required"`
	Height  uint64      `json:"height"`
	TxIndex uint16      `json:"tx_index"`
}

func NewReceipt(root []byte, status ReceiptStatus, cumulativeGasUsed uint64) *Receipt {
	r := &Receipt{PostState: common.CopyBytes(root), CumulativeGasUsed: cumulativeGasUsed, Status: status}
	return r
}

func (r *Receipt) Size() common.StorageSize {
	size := common.StorageSize(unsafe.Sizeof(*r)) + common.StorageSize(len(r.PostState))
	return size
}

func (r *Receipt) Success() bool {
	return r.Status == RSSuccess
}

type Receipts []*Receipt

func (r Receipts) Len() int { return len(r) }
