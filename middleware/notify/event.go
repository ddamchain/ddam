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

package notify

// defines all of current used event ids
const (
	BlockAddSucc = "block_add_succ"

	NewBlock = "new_block"

	BlockReq = "block_req"

	BlockResponse = "block_response"

	BlockInfoNotify = "block_info_notify"

	ChainPieceBlockReq = "chain_piece_block_req"

	ChainPieceBlock = "chain_piece_block"

	BlockSync = "block_sync"

	TxSyncNotify   = "tx_sync_notify"
	TxSyncReq      = "tx_sync_req"
	TxSyncResponse = "tx_sync_response"

	ConsoleMessage = "console_msg"
)
