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

package network

const (

	//The following four messages are used for block sync
	BlockInfoNotifyMsg uint32 = 1
	ReqBlock           uint32 = 2
	BlockResponseMsg   uint32 = 3
	NewBlockMsg        uint32 = 4

	//The following two messages are used for block fork processing
	ReqChainPieceBlock uint32 = 5
	ChainPieceBlock    uint32 = 6

	//The following three message are used for tx sync
	TxSyncNotify   uint32 = 7
	TxSyncReq      uint32 = 8
	TxSyncResponse uint32 = 9
)

type Message struct {
	ChainID uint16

	ProtocolVersion uint16

	Code uint32

	Body []byte
}

type Conn struct {
	ID   string
	IP   string
	Port string
}

type MsgDigest []byte

type MsgHandler interface {
	Handle(sourceID string, msg Message) error
}

type Network interface {
	//Send message to the node which id represents.If self doesn't connect to the node,
	// resolve the kad net to find the node and then send the message
	Send(id string, msg Message) error

	//TransmitToNeighbor Send message to neighbor nodes
	TransmitToNeighbor(msg Message) error

	//Broadcast Send the message to all nodes it connects to and the node which receive the message also broadcast the message to their neighbor once
	Broadcast(msg Message) error

	//ConnInfo Return all connections self has
	ConnInfo() []Conn
}
