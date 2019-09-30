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

package consensus

import (
	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/consensus/plotter"
	"github.com/xchain/go-chain/global/types"
	"github.com/xchain/go-chain/middleware/time"
	"github.com/xchain/go-chain/network"
	"math/big"
)

type plotManagerInterface interface {
	// LoadScoops load scoop from plotter file of the given params
	LoadScoops(scoopNum int, startNonce, endNonce types.Nonce) []*plotter.ScoopFragment

	// GenerateScoop generates scoop for the given params
	GenerateScoop(address common.Address, nonce types.Nonce, scoopNum int) *plotter.Scoop

	// NonceRange returns the nonce range for the plotters the node loads
	NonceRange() (min, max types.Nonce)
}

type chainAccessor interface {
	types.LatestDBGetter
	BatchGetBlockHeadersAfterHeight(h uint64, limit int) []*types.BlockHeader
	BlockProposal(nonce types.Nonce, curTime time.TimeStamp, baseTarget types.Target, cumulativeDiff *big.Int, code types.AuthCode, address common.Address) *types.Block
	AddBlockOnChain(source string, b *types.Block) types.AddBlockResult
	RegisterEngine(engine types.ConsensusEngine)
	QueryTopBlock() *types.BlockHeader
}

type authorityI interface {
	// CanPropose checks if the address -> sha256(address+umid) pair on the chain, if so then return true, false otherwise
	// umid is generated real-time in current machine
	CanPropose(address common.Address) bool

	// Generate umid mixer when propose block by the given block hash
	// returns sha256(sha256(address+umid)+blockHash)
	// umid is generated real-time in current machine
	GenerateAuthCode(blockHash common.Hash, address common.Address) types.AuthCode

	// Verify the mixer when insert block
	// First get the sha256(address+umid) from chain, and then compares if mixer == sha256(sha256(address+umid)+blockHash)
	// If equality, then mixer is legal, otherwise illegal
	VerifyAuthCode(ac types.AuthCode, address common.Address, blockHash common.Hash) (bool, error)

	RegisterGetter(getter types.LatestDBGetter)
}

type broadcastI interface {
	Broadcast(msg network.Message) error
}

type engineMgr interface {
	finishCh() chan *worker
	findCh() chan *worker
	LoadScoops(scoopNum int, startNonce, endNonce types.Nonce) []*plotter.ScoopFragment
}
