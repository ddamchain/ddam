//   Copyright (C) 2019 XChain
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

package core

import (
	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/global/types"
)

type forkProcessI interface {
	process(string, *types.Block)
}

type chainAccessor interface {
	types.ChainReader
	AddBlockOnChain(source string, b *types.Block) types.AddBlockResult
	BatchGetBlocksAfterHeight(h uint64, limit int) []*types.Block
	AddChainSliceOnChain(source string, chainSlice []*types.Block, weight types.ChainWeight) []*types.Block
	IsAdjusting() bool
	SyncFinished() bool
	verifyBlockHeader(bh *types.BlockHeader) error
}

type blockSyncI interface {
	getPeerTopBlock(id string) *types.BlockHeader
	trySyncRoutine() bool
}

type txPoolI interface {
	// PackForCast returns a list of transactions for casting a block
	PackForCast() []*types.Transaction

	// GetTransaction trys to find a transaction from pool by hash and return it
	GetTransaction(hash common.Hash) *types.Transaction

	GetReceipt(hash common.Hash) *types.Receipt

	// RemoveFromPool removes the transactions from pool by hash
	RemoveFromPool(txs []common.Hash)

	// BackToPool will put the transactions back to pool
	BackToPool(txs []*types.Transaction)

	// RecoverAndValidateTx recovers the sender of the transaction and also validates the transaction
	RecoverAndValidateTx(tx *types.Transaction) error

	SaveReceipts(blockHash common.Hash, receipts types.Receipts) error

	DeleteReceipts(txs []common.Hash) error
}

type umidStore interface {
	// BindUMIDAddress stores the address -> sha256(address+umid) to the chain
	BindUMIDAddress(db types.AccountDB, address common.Address, umid types.UMID) error

	// TransferUMIDAddress replaces the current umid owner
	TransferUMIDAddress(db types.AccountDB, from, to common.Address, umid types.UMID) error

	//UnbindUMIDAddress remove the sha256(address+umid) from the chain
	UnbindUMIDAddress(db types.AccountDB, from common.Address, umid types.UMID) error
}

type txGetter interface {
	GetTransactionByHash(needSource bool, h common.Hash) *types.Transaction
}

type pledgeMgr interface {
	getPledge(db types.AccountDB, address common.Address) uint64
}
