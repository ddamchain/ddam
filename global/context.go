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

package global

import (
	"fmt"
	"github.com/xchain/go-chain/global/types"
	"github.com/xchain/go-chain/middleware/notify"
	"github.com/xchain/go-chain/middleware/ticker"
	"github.com/xchain/go-chain/middleware/time"
	"reflect"
)

var ctx *ContextFactory

func Context() *ContextFactory {
	return ctx
}

type ContextFactory struct {
	Config      types.ConfManager
	Current     *types.Miner
	TimeService time.TimeService
	Ticker      *ticker.GlobalTicker
	Bus         *notify.Bus
}

func Init(confFile string) {
	ctx = &ContextFactory{
		Config:      types.NewConfINIManager(confFile),
		Bus:         notify.NewBus(),
		TimeService: time.InitTimeSync(),
		Ticker:      ticker.NewGlobalTicker("global"),
	}
}

func (ctx *ContextFactory) Register(field string, value interface{}) {
	ty := reflect.ValueOf(ctx).Elem()
	v := ty.FieldByName(field)
	if v.CanSet() {
		v.Set(reflect.ValueOf(value))
	} else {
		panic(fmt.Errorf("field %v cannot be set", field))
	}
}
