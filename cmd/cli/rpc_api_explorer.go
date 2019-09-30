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
	"fmt"
	"github.com/vmihailenco/msgpack"
	"github.com/xchain/go-chain/auth"
	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/core"
	"github.com/xchain/go-chain/global/types"
	"math/big"
)

var KeyOfUMIDAddr = []byte{'U', 'M', 'I', 'D'}

// RpcExplorerImpl provides rpc service for blockchain explorer use
type RpcExplorerImpl struct {
	*baseRpcImpl
}

func (api *RpcExplorerImpl) Namespace() string {
	return "Explorer"
}

func (api *RpcExplorerImpl) Version() string {
	return "1"
}

// ExplorerBlockDetail is used in the blockchain browser to query block details
func (api *RpcExplorerImpl) ExplorerBlockDetail(height uint64) (*Result, error) {
	b := api.br.QueryBlockByHeight(height)
	if b == nil {
		return failResult("QueryBlock error")
	}
	block := convertBlockHeader(b)

	trans := make([]Transaction, 0)

	for _, tx := range b.Transactions {
		trans = append(trans, *convertTransaction(tx))
	}

	evictedReceipts := make([]*types.Receipt, 0)

	receipts := make([]*types.Receipt, len(b.Transactions))
	for i, tx := range b.Transactions {
		wrapper := api.txPool.GetReceipt(tx.Hash)
		if wrapper != nil {
			receipts[i] = wrapper
		}
	}

	bd := &ExplorerBlockDetail{
		BlockDetail:     BlockDetail{Block: *block, Trans: trans},
		EvictedReceipts: evictedReceipts,
		Receipts:        receipts,
	}
	return successResult(bd)
}

func (api *RpcExplorerImpl) GetUmidAddresses(addr common.Address) (*Result, error) {
	db := api.br.LatestStateDB()

	umidAddreses := dbGet(db, addr)
	if umidAddreses == nil {
		return failResult(fmt.Errorf("No UMID registration information ").Error())
	}
	target, err := unMarshalUmidAddreses(umidAddreses)
	if err != nil {
		return failResult(err.Error())
	}
	return successResult(target)
}

func (api *RpcExplorerImpl) GetPledge(addr common.Address) (*Result, error) {
	db := api.br.LatestStateDB()

	data := db.GetData(addr, core.KeyOfstakeAmount)
	if data == nil {
		return failResult("Pledged information not found !!")
	}
	mount := new(big.Int).SetBytes(data)
	return successResult(mount)
}

func dbGet(db types.AccountDB, address common.Address) []byte {
	return db.GetData(address, KeyOfUMIDAddr)
}

func unMarshalUmidAddreses(umidAddreses []byte) ([][]byte, error) {
	var target auth.UMIDAddreses
	err := msgpack.Unmarshal(umidAddreses, &target)
	if err != nil {
		return nil, err
	}
	return target, nil
}
