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
	"github.com/xchain/go-chain/crypto"
	"math/big"
	"time"

	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/global/types"
)

// Result is rpc request successfully returns the variable parameter
type Result struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

func (r *Result) IsSuccess() bool {
	return r.Status == 0
}

// ErrorResult is rpc request error returned variable parameter
type ErrorResult struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// RPCReqObj is complete rpc request body
type RPCReqObj struct {
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Jsonrpc string        `json:"jsonrpc"`
	ID      uint          `json:"id"`
}

// RPCResObj is complete rpc response body
type RPCResObj struct {
	Jsonrpc string       `json:"jsonrpc"`
	ID      uint         `json:"id"`
	Result  *Result      `json:"result,omitempty"`
	Error   *ErrorResult `json:"error,omitempty"`
}

// Transactions in the buffer pool transaction list
type Transactions struct {
	Hash      string `json:"hash"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	Value     string `json:"value"`
	Height    uint64 `json:"height"`
	BlockHash string `json:"block_hash"`
}

type PubKeyInfo struct {
	PubKey string `json:"pub_key"`
	ID     string `json:"id"`
}

type ConnInfo struct {
	ID      string `json:"id"`
	IP      string `json:"ip"`
	TCPPort string `json:"tcp_port"`
}

type GroupStat struct {
	Dismissed bool  `json:"dismissed"`
	VCount    int32 `json:"v_count"`
}

type ProposerStat struct {
	Stake      uint64  `json:"stake"`
	StakeRatio float64 `json:"stake_ratio"`
	PCount     int32   `json:"p_count"`
}

type NodeInfo struct {
	ID           string  `json:"id"`
	Balance      float64 `json:"balance"`
	Status       string  `json:"status"`
	WGroupNum    int     `json:"w_group_num"`
	AGroupNum    int     `json:"a_group_num"`
	NType        string  `json:"n_type"`
	TxPoolNum    int     `json:"tx_pool_num"`
	BlockHeight  uint64  `json:"block_height"`
	GroupHeight  uint64  `json:"group_height"`
	VrfThreshold float64 `json:"vrf_threshold"`
}

type PageObjects struct {
	Total uint64        `json:"count"`
	Data  []interface{} `json:"data"`
}

type Block struct {
	Height               uint64         `json:"height"`
	Hash                 common.Hash    `json:"hash"`
	PreHash              common.Hash    `json:"pre_hash"`
	CurTime              time.Time      `json:"cur_time"`
	Proposer             common.Address `json:"proposer"`
	Nonce                types.Nonce    `json:"nonce"`
	TxTree               common.Hash    `json:"tx_tree"`
	ReceiptTree          common.Hash    `json:"receipt_tree"`
	StateTree            common.Hash    `json:"state_tree"`
	BaseTarget           types.Target   `json:"base_target"`
	CumulativeDifficulty *big.Int       `json:"cumulative_difficulty"`
	Difficulty           *big.Int       `json:"difficulty"`
	Capacity             string         `json:"capacity"`
	Auth                 types.AuthCode `json:"auth"`
	Sign                 *crypto.Sign   `json:"sign"`
}

type BlockDetail struct {
	Block
	Trans      []Transaction `json:"trans"`
	PreTotalQN uint64        `json:"pre_total_qn"`
}

type BlockReceipt struct {
	Receipts        []*types.Receipt `json:"receipts"`
	EvictedReceipts []*types.Receipt `json:"evictedReceipts"`
}

type ExplorerBlockDetail struct {
	BlockDetail
	Receipts        []*types.Receipt `json:"receipts"`
	EvictedReceipts []*types.Receipt `json:"evictedReceipts"`
}

type Group struct {
	Seed          common.Hash `json:"id"`
	BeginHeight   uint64      `json:"begin_height"`
	DismissHeight uint64      `json:"dismiss_height"`
	Threshold     int32       `json:"threshold"`
	Members       []string    `json:"members"`
	MemSize       int         `json:"mem_size"`
}

type Transaction struct {
	Data   []byte          `json:"data"`
	Value  float64         `json:"value"`
	Nonce  uint64          `json:"nonce"`
	Source *common.Address `json:"source"`
	Target *common.Address `json:"target"`
	Type   int8            `json:"type"`

	GasLimit uint64      `json:"gas_limit"`
	GasPrice uint64      `json:"gas_price"`
	Hash     common.Hash `json:"hash"`
}

type Receipt struct {
	Status            int    `json:"status"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed"`

	TxHash  common.Hash `json:"transactionHash" gencodec:"required"`
	Height  uint64      `json:"height"`
	TxIndex uint16      `json:"tx_index"`
}

type ExecutedTransaction struct {
	Receipt     *Receipt
	Transaction *Transaction
}

type Dashboard struct {
	BlockHeight uint64     `json:"block_height"`
	GroupHeight uint64     `json:"group_height"`
	WorkGNum    int        `json:"work_g_num"`
	NodeInfo    *NodeInfo  `json:"node_info"`
	Conns       []ConnInfo `json:"conns"`
}

type ExplorerAccount struct {
	Balance   *big.Int               `json:"balance"`
	Nonce     uint64                 `json:"nonce"`
	Type      uint32                 `json:"type"`
	CodeHash  string                 `json:"code_hash"`
	Code      string                 `json:"code"`
	StateData map[string]interface{} `json:"state_data"`
}
