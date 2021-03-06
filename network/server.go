//   Copyright (C) 2018 XChain
//
//   This program is free software: you cas redistribute it and/or modify
//   it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either versios 3 of the License, or
//   (at your option) any later versios.
//
//   This program is distributed is the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without eves the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//   GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
//   along with this program.  If not, see <https://www.gnu.org/licenses/>.

package network

import (
	"github.com/golang/protobuf/proto"
	"github.com/xchain/go-chain/global"
	"github.com/xchain/go-chain/global/types"
	"github.com/xchain/go-chain/middleware/pb"

	"strconv"
	"time"

	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/middleware/notify"

	"golang.org/x/crypto/sha3"
)

type Server struct {
	Self *Node

	netCore *NetCore
	config  *NetworkConfig
}

func (s *Server) Send(id string, msg Message) error {
	bytes, err := marshalMessage(msg)
	if err != nil {
		return err
	}
	if id == s.Self.ID.GetHexString() {
		s.sendSelf(bytes)
		return nil
	}
	go s.netCore.sendToNode(NewNodeID(id), nil, bytes, msg.Code)

	return nil
}

func (s *Server) TransmitToNeighbor(msg Message) error {
	bytes, err := marshalMessage(msg)
	if err != nil {
		Logger.Errorf("Marshal message error:%s", err.Error())
		return err
	}

	s.netCore.broadcast(bytes, msg.Code, false, nil, -1)

	return nil
}

func (s *Server) Broadcast(msg Message) error {
	bytes, err := marshalMessage(msg)
	if err != nil {
		Logger.Errorf("Marshal message error:%s", err.Error())
		return err
	}
	s.netCore.broadcast(bytes, msg.Code, true, nil, -1)

	return nil
}

func (s *Server) ConnInfo() []Conn {
	result := make([]Conn, 0)
	peers := s.netCore.peerManager.peers
	for _, p := range peers {
		if p.sessionID > 0 && p.IP != nil && p.Port > 0 && p.isAuthSucceed {
			c := Conn{ID: p.ID.GetHexString(), IP: p.IP.String(), Port: strconv.Itoa(p.Port)}
			result = append(result, c)
		}
	}
	return result
}

func (s *Server) sendSelf(b []byte) {
	s.handleMessage(b, s.Self.ID.GetHexString(), s.netCore.chainID, s.netCore.protocolVersion)
}

func (s *Server) handleMessage(b []byte, from string, chainID uint16, protocolVersion uint16) {

	message, error := unMarshalMessage(b)
	if error != nil {
		Logger.Errorf("Proto unmarshal error:%s", error.Error())
		return
	}
	message.ChainID = chainID
	message.ProtocolVersion = protocolVersion
	Logger.Debugf("Receive message from %s,code:%d,msg size:%d,hash:%s, chainID:%v,protocolVersion:%v", from, message.Code, len(b), message.Hash(), chainID, protocolVersion)
	s.netCore.flowMeter.recv(int64(message.Code), int64(len(b)))

	go s.handleMessageInner(message, from)
}

func newNotifyMessage(message *Message, from string) *types.DefaultMessage {
	return types.NewDefaultMessage(message.Body, from, message.ChainID, message.ProtocolVersion)
}

func (s *Server) handleMessageInner(message *Message, from string) {

	s.netCore.onHandleDataMessageStart()
	defer s.netCore.onHandleDataMessageDone()

	begin := time.Now()
	code := message.Code

	topicID := ""
	switch code {
	case TxSyncNotify:
		topicID = notify.TxSyncNotify
	case TxSyncReq:
		topicID = notify.TxSyncReq
	case TxSyncResponse:
		topicID = notify.TxSyncResponse
	case BlockInfoNotifyMsg:
		topicID = notify.BlockInfoNotify
	case ReqBlock:
		topicID = notify.BlockReq
	case BlockResponseMsg:
		topicID = notify.BlockResponse
	case NewBlockMsg:
		topicID = notify.NewBlock
	case ReqChainPieceBlock:
		topicID = notify.ChainPieceBlockReq
	case ChainPieceBlock:
		topicID = notify.ChainPieceBlock
	}
	if topicID != "" {
		msg := newNotifyMessage(message, from)
		global.Context().Bus.Publish(topicID, msg)
	}

	if time.Since(begin) > 100*time.Millisecond {
		Logger.Debugf("handle message cost time:%v,hash:%s,code:%d", time.Since(begin), message.Hash(), code)
	}
}

func marshalMessage(m Message) ([]byte, error) {
	message := pb.Message{Code: &m.Code, Body: m.Body}
	return proto.Marshal(&message)
}

func unMarshalMessage(b []byte) (*Message, error) {
	message := new(pb.Message)
	e := proto.Unmarshal(b, message)
	if e != nil {
		return nil, e
	}
	m := Message{Code: *message.Code, Body: message.Body}
	return &m, nil
}

func (m Message) Hash() string {
	bytes, err := marshalMessage(m)
	if err != nil {
		return ""
	}

	var h common.Hash
	sha3Hash := sha3.Sum256(bytes)
	copy(h[:], sha3Hash[:])

	return h.Hex()
}
