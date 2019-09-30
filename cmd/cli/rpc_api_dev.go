//   Copyright (C) 2019 XChain
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
	"github.com/xchain/go-chain/global"
	"strings"

	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/core"
	"github.com/xchain/go-chain/global/types"
	"github.com/xchain/go-chain/network"
)

// RpcDevImpl provides api functions for those develop core features.
// It is mainly for debug or test use
type RpcDevImpl struct {
	*baseRpcImpl
}

func (api *RpcDevImpl) Namespace() string {
	return "Dev"
}

func (api *RpcDevImpl) Version() string {
	return "1"
}

// ConnectedNodes query the information of the connected node
func (api *RpcDevImpl) ConnectedNodes() (*Result, error) {

	nodes := network.GetNetInstance().ConnInfo()
	conns := make([]ConnInfo, 0)
	for _, n := range nodes {
		conns = append(conns, ConnInfo{ID: n.ID, IP: n.IP, TCPPort: n.Port})
	}
	return successResult(conns)
}

// get transaction by hash
func (api *RpcDevImpl) GetTransaction(hash string) (*Result, error) {
	if !validateHash(strings.TrimSpace(hash)) {
		return failResult("Wrong hash format")
	}
	transaction := api.br.GetTransactionByHash(true, common.HexToHash(hash))
	if transaction == nil {
		return failResult("transaction not exists")
	}
	detail := make(map[string]interface{})
	detail["hash"] = hash
	if transaction.Source != nil {
		detail["source"] = transaction.Source.Hash().Hex()
	}
	if transaction.Target != nil {
		detail["target"] = transaction.Target.Hash().Hex()
	}
	detail["value"] = transaction.Value

	return successResult(detail)
}

func (api *RpcDevImpl) GetBlocks(from uint64, to uint64) (*Result, error) {
	if from > to {
		return failResult("param error")
	}
	blocks := make([]*Block, 0)
	var preBH *types.BlockHeader
	for h := from; h <= to; h++ {
		b := api.br.QueryBlockByHeight(h)
		if b != nil {
			block := convertBlockHeader(b)
			if preBH == nil {
				preBH = api.br.QueryBlockHeaderByHash(b.Header.PreHash)
			}
			preBH = b.Header
			blocks = append(blocks, block)
		}
	}
	return successResult(blocks)
}

func (api *RpcDevImpl) GetTopBlock() (*Result, error) {
	bh := api.br.QueryTopBlock()
	b := api.br.QueryBlockByHash(bh.Hash)
	bh = b.Header

	blockDetail := make(map[string]interface{})
	blockDetail["hash"] = common.ShortHex(bh.Hash.Hex())
	blockDetail["height"] = bh.Height
	blockDetail["pre_hash"] = common.ShortHex(bh.PreHash.Hex())
	blockDetail["cur_time"] = bh.CurTime.Local().Format("2006-01-02 15:04:05")
	blockDetail["castor"] = common.ShortHex(bh.Proposer.Hex())
	blockDetail["txs"] = len(b.Transactions)
	blockDetail["target"] = bh.BaseTarget
	blockDetail["capacity"] = core.EstimateCapacity(bh.BaseTarget).String()

	blockDetail["tx_pool_total"] = api.txPool.TxNum()
	blockDetail["miner_id"] = common.ShortHex(global.Context().Current.Addr.Hex())
	return successResult(blockDetail)
}

func (api *RpcDevImpl) BlockReceipts(h string) (*Result, error) {
	if !validateHash(strings.TrimSpace(h)) {
		return failResult("Wrong param format")
	}
	chain := core.BlockChainImpl
	b := chain.QueryBlockByHash(common.HexToHash(h))
	if b == nil {
		return failResult("block not found")
	}

	evictedReceipts := make([]*types.Receipt, 0)
	receipts := make([]*types.Receipt, len(b.Transactions))
	for i, tx := range b.Transactions {
		wrapper := api.txPool.GetReceipt(tx.Hash)
		if wrapper != nil {
			receipts[i] = wrapper
		}
	}
	br := &BlockReceipt{EvictedReceipts: evictedReceipts, Receipts: receipts}
	return successResult(br)
}
