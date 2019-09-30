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
	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/global"
	"io"
	"os"

	"github.com/xchain/go-chain/core"
	"github.com/xchain/go-chain/global/types"
	"github.com/xchain/go-chain/middleware/notify"
)

type msgShower struct {
	out    io.Writer
	bchain types.ChainReader
}

var shower *msgShower

func showMsg(format string, a ...interface{}) {
	if shower == nil {
		fmt.Fprintf(os.Stdout, format+"\n", a...)
	} else {
		shower.showMsg(format+"\n", a...)
	}
}

func initMsgShower(cr types.ChainReader) {
	ii := &msgShower{
		out:    os.Stdout,
		bchain: cr,
	}

	global.Context().Bus.Subscribe(notify.BlockAddSucc, ii.onBlockAddSuccess)
	global.Context().Bus.Subscribe(notify.BlockSync, ii.blockSync)
	global.Context().Bus.Subscribe(notify.ConsoleMessage, ii.consoleMsg)

	shower = ii

	shower.showMsg("local top %v %v", cr.QueryTopBlock().Hash, cr.QueryTopBlock().Height)
}

func (ms *msgShower) showMsg(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	fmt.Fprintf(ms.out, fmt.Sprintf("%v\n", s))
}

func (ms *msgShower) onBlockAddSuccess(message notify.Message) {
	b := message.GetData().(*types.Block)
	if b.Header.Proposer.EqualTo(global.Context().Current.Addr) {
		ms.showMsg("congratulations, you mined block height %v success!", b.Header.Height)
	}
	ms.showMsg("add block success: %v %v, castor: %v", b.Header.Hash, b.Header.Height, common.ShortHex(b.Header.Proposer.Hex()))
}

func (ms *msgShower) blockSync(message notify.Message) {
	cand := message.GetData().(*core.SyncCandidateInfo)
	ms.showMsg("sync block from %v[height=%v], localHeight=%v, reqHeight %v", cand.Candidate, cand.CandidateHeight, ms.bchain.Height(), cand.ReqHeight)
}

func (ms *msgShower) consoleMsg(message notify.Message) {
	msg := message.GetData().(string)
	ms.showMsg(msg)
}
