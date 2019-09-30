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
	"math/big"

	"github.com/xchain/go-chain/common"
)

// AddBlockResult is the result of the add-block operation
type AddBlockResult int8

// defines all possible result of the add-block operation
const (
	AddBlockFailed  AddBlockResult = 0 // Means the operations is fail
	AddBlockSuccess AddBlockResult = 1 // Means success
	BlockExisted    AddBlockResult = 2 // Means the block already added before
)

type ChainWeight interface {
	Weight() *big.Int
	MoreWeightThan(chain ChainWeight) bool
}

type ChainReader interface {
	Height() uint64
	QueryTopBlock() *BlockHeader
	QueryBlockHeaderByHash(hash common.Hash) *BlockHeader
	QueryBlockHeaderByHeight(height uint64) *BlockHeader
	HasBlock(hash common.Hash) bool
	HasHeight(height uint64) bool
}

type LatestDBGetter interface {
	LatestStateDB() AccountDB
}

// BlockChain is a interface, encapsulates some methods for manipulating the blockchain
//type BlockChain interface {
//	ChainReader
//	AccountRepository
//
//	// LatestStateDB returns core's last account database
//	LatestStateDB() AccountDB
//
//	// QueryBlockByHash query the block by hash
//	QueryBlockByHash(hash common.Hash) *Block
//
//	// QueryBlockByHeight query the block by height
//	QueryBlockByHeight(height uint64) *Block
//
//	// BatchGetBlocksAfterHeight query blocks after the specified height
//	BatchGetBlocksAfterHeight(height uint64, limit int) []*Block
//
//	// BatchGetBlocksAfterHeight query blocks after the specified height
//	BatchGetBlockHeadersAfterHeight(height uint64, limit int) []*BlockHeader
//
//	// GetTransactionByHash get a transaction by hash
//	GetTransactionByHash(needSource bool, h common.Hash) *Transaction
//
//	// GetTransactionPool return the transaction pool waiting for the block
//	GetTransactionPool() TransactionPool
//
//	// IsAdjusting means whether need to adjust blockchain, which means there may be a fork
//	IsAdjusting() bool
//
//	// Remove removes the block and blocks after it from the core. Only used in a debug file, should be removed later
//	Remove(block *Block) bool
//
//	// Close the open levelDb files
//	Close()
//
//	// StateDBAt returns account database with specified block height
//	StateDBAt(height uint64) (AccountDB, error)
//
//	// Version of core Id
//	Version() int
//}

// ExecutedTransaction contains the transaction and its receipt
type ExecutedTransaction struct {
	Receipt     *Receipt
	Transaction *Transaction
}

// AccountRepository contains account query interface
type AccountRepository interface {
	// GetBalance return the balance of specified address
	GetBalance(address common.Address) *big.Int

	// GetBalance returns the nonce of specified address
	GetNonce(address common.Address) uint64
}

//
//
//type txPoolI interface {
//	// PackForCast returns a list of transactions for casting a block
//	PackForCast() []*types.Transaction
//
//	// GetTransaction trys to find a transaction from pool by hash and return it
//	GetTransaction(hash common.Hash) *types.Transaction
//
//	// GetReceipt returns the transaction's recipe by hash
//	GetReceipt(hash common.Hash) *types.Receipt
//
//	// GetReceived returns the received transactions in the pool with a limited size
//	GetReceived() []*types.Transaction
//
//	// TxNum returns the number of transactions in the pool
//	TxNum() uint64
//
//	// RemoveFromPool removes the transactions from pool by hash
//	RemoveFromPool(txs []common.Hash)
//
//	// BackToPool will put the transactions back to pool
//	BackToPool(txs []*Transaction)
//
//	// RecoverAndValidateTx recovers the sender of the transaction and also validates the transaction
//	RecoverAndValidateTx(tx *Transaction) error
//
//	SaveReceipts(blockHash common.Hash, receipts Receipts) error
//
//	DeleteReceipts(txs []common.Hash) error
//
//	//check transaction hash exist in local
//	IsTransactionExisted(hash common.Hash) (exists bool, where int)
//}
