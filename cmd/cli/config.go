//   Copyright (C) 2019 XChain
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

type minerConfig struct {
	confFile      string
	rpcLevel      rpcLevel
	rpcAddr       string
	rpcPort       uint16
	super         bool
	testMode      bool
	natIP         string
	natPort       uint16
	seedIP        string
	seedID        string
	keystore      string
	enableMonitor bool
	chainID       uint16
	password      string
}

func (cfg *minerConfig) rpcEnable() bool {
	return cfg.rpcLevel > rpcLevelNone
}
