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
	"encoding/json"
	"fmt"
	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/core"
	"github.com/xchain/go-chain/global/types"
	"math/big"
	"strings"
)

type blockReader interface {
	types.ChainReader
	types.AccountRepository
	types.LatestDBGetter

	QueryBlockByHeight(height uint64) *types.Block
	QueryBlockByHash(hash common.Hash) *types.Block
	GetTransactionByHash(needSource bool, h common.Hash) *types.Transaction
	GetAccountDBByHash(hash common.Hash) (types.AccountDB, error)
}

type txPool interface {
	TxNum() uint64
	GetReceipt(hash common.Hash) *types.Receipt
	AddTransaction(tx *types.Transaction) (bool, error)
}

type baseRpcImpl struct {
	br     blockReader
	txPool txPool
}

// RpcGtasImpl provides rpc service for users to interact with remote nodes
type RpcGxImpl struct {
	*baseRpcImpl
}

func (api *RpcGxImpl) Namespace() string {
	return "Gx"
}

func (api *RpcGxImpl) Version() string {
	return "1"
}

func successResult(data interface{}) (*Result, error) {
	return &Result{
		Message: "success",
		Data:    data,
		Status:  0,
	}, nil
}
func failResult(err string) (*Result, error) {
	return &Result{
		Message: err,
		Data:    nil,
		Status:  -1,
	}, nil
}

// Tx is user transaction interface, used for sending transaction to the node
func (api *RpcGxImpl) Tx(txRawjson string) (*Result, error) {
	var txRaw = new(txRawData)
	if err := json.Unmarshal([]byte(txRawjson), txRaw); err != nil {
		return failResult(err.Error())
	}
	if !validateTxType(txRaw.TxType) {
		return failResult("Not supported txType")
	}

	// Check the address for the specified tx types
	switch txRaw.TxType {
	case types.TransactionTypeTransfer, types.TransactionTypeBindUMID, types.TransactionTypeTransformUMID, types.TransactionTypeUnbindUMID, types.TransactionTypeStakeAdd, types.TransactionTypeStakeReduce:
		if !validateAddress(strings.TrimSpace(txRaw.Target)) {
			return failResult("Wrong target address format")
		}
	}

	trans := txRawToTransaction(txRaw)

	trans.Hash = trans.GenHash()

	if err := api.sendTransaction(trans); err != nil {
		return failResult(err.Error())
	}

	return successResult(trans.Hash.Hex())
}

// Balance is query balance interface
func (api *RpcGxImpl) Balance(account string) (*Result, error) {
	if !validateAddress(strings.TrimSpace(account)) {
		return failResult("Wrong account address format")
	}
	b := api.br.GetBalance(common.HexToAddress(account))

	balance := common.AM2DDAM(b.Uint64())
	return &Result{
		Message: fmt.Sprintf("The balance of account: %s is %v DDAM", account, balance),
		Data:    balance,
	}, nil
}

// BlockHeight query block height
func (api *RpcGxImpl) BlockHeight() (*Result, error) {
	height := api.br.QueryTopBlock().Height
	return successResult(height)
}

func (api *RpcGxImpl) GetBlockByHeight(height uint64) (*Result, error) {
	b := api.br.QueryBlockByHeight(height)
	if b == nil {
		return failResult("height not exists")
	}

	block := convertBlockHeader(b)

	return successResult(block)
}

func (api *RpcGxImpl) GetBlockByHash(hash string) (*Result, error) {
	if !validateHash(strings.TrimSpace(hash)) {
		return failResult("Wrong hash format")
	}
	b := api.br.QueryBlockByHash(common.HexToHash(hash))
	if b == nil {
		return failResult("height not exists")
	}

	block := convertBlockHeader(b)

	return successResult(block)
}

func (api *RpcGxImpl) TransDetail(h string) (*Result, error) {
	if !validateHash(strings.TrimSpace(h)) {
		return failResult("Wrong hash format")
	}
	tx := api.br.GetTransactionByHash(true, common.HexToHash(h))

	if tx != nil {
		trans := convertTransaction(tx)
		return successResult(trans)
	}
	return successResult(nil)
}

func (api *RpcGxImpl) Nonce(addr string) (*Result, error) {
	if !validateAddress(strings.TrimSpace(addr)) {
		return failResult("Wrong account address format")
	}
	address := common.HexToAddress(addr)
	// user will see the nonce as db nonce +1, so that user can use it directly when send a transaction
	nonce := api.br.GetNonce(address) + 1
	return successResult(nonce)
}

func (api *RpcGxImpl) TxReceipt(h string) (*Result, error) {
	if !validateHash(strings.TrimSpace(h)) {
		return failResult("Wrong hash format")
	}
	hash := common.HexToHash(h)
	rc := api.txPool.GetReceipt(hash)
	if rc != nil {
		tx := api.br.GetTransactionByHash(true, hash)
		return successResult(convertExecutedTransaction(&types.ExecutedTransaction{
			Receipt:     rc,
			Transaction: tx,
		}))
	}
	return failResult("tx not exist")
}

func (api *baseRpcImpl) sendTransaction(trans *types.Transaction) error {
	if trans.Sign == nil {
		return fmt.Errorf("transaction sign is empty")
	}

	if ok, err := api.txPool.AddTransaction(trans); err != nil || !ok {
		//common.DefaultLogger.Errorf("AddTransaction not ok or error:%s", err.Error())，
		return err
	}
	return nil
}

func (api *baseRpcImpl) Stake(account string) (*Result, error) {
	if !validateAddress(strings.TrimSpace(account)) {
		return failResult("Wrong account address format")
	}
	db := api.br.LatestStateDB()

	data := db.GetData(common.HexToAddress(account), core.KeyOfstakeAmount)
	if data == nil {
		return failResult("Stake information not found !!")
	}
	mount := new(big.Int).SetBytes(data)

	stake := common.AM2DDAM(mount.Uint64())

	return &Result{
		Message: fmt.Sprintf("The stake of account: %s is %v DDAM", account, stake),
		Data:    stake,
	}, nil

	//return successResult(mount)
}
