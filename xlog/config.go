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

package xlog

const (
	DefaultConfig = `<seelog minlevel="debug">
						<outputs formatid="default">
							<rollingfile type="size" filename="./logs/default.log" maxsize="100000000" maxrolls="1"/>
						</outputs>
						<formats>
							<format id="default" format="%Date(2006-01-02 15:04:05.000)[%File:%Line] %Msg%n" />
						</formats>
					</seelog>`

	CoreLogConfig = `<seelog minlevel="debug">
						<outputs formatid="default">
							<rollingfile type="size" filename="./logs/core.log" maxsize="200000000" maxrolls="1"/>
						</outputs>
						<formats>
							<format id="default" format="%Date(2006-01-02 15:04:05.000)[%File:%Line]%Msg%n" />
						</formats>
					</seelog>`

	BlockSyncLogConfig = `<seelog minlevel="debug">
						<outputs formatid="default">
							<rollingfile type="size" filename="./logs/block_sync.log" maxsize="200000000" maxrolls="1"/>
						</outputs>
						<formats>
							<format id="default" format="%Date(2006-01-02 15:04:05.000)[%File:%Line]%Msg%n" />
						</formats>
					</seelog>`

	P2PLogConfig = `<seelog minlevel="debug">
						<outputs formatid="default">
							<rollingfile type="size" filename="./logs/p2p.log" maxsize="200000000" maxrolls="1"/>
						</outputs>
						<formats>
							<format id="default" format="%Date(2006-01-02 15:04:05.000)[%File:%Line]%Msg%n" />
						</formats>
					</seelog>`

	EngineLogConfig = `<seelog minlevel="debug">
						<outputs formatid="default">
							<rollingfile type="size" filename="./logs/engine.log" maxsize="200000000" maxrolls="1"/>
						</outputs>
						<formats>
							<format id="default" format="%Date(2006-01-02 15:04:05.000)[%File:%Line]%Msg%n" />
						</formats>
					</seelog>`

	ForkLogConfig = `<seelog minlevel="debug">
						<outputs formatid="default">
							<rollingfile type="size" filename="./logs/fork.log" maxsize="100000000" maxrolls="1"/>
						</outputs>
						<formats>
							<format id="default" format="%Date(2006-01-02 15:04:05.000)[%File:%Line]%Msg%n" />
						</formats>
					</seelog>`
	TxSyncLogConfig = `<seelog minlevel="debug">
						<outputs formatid="default">
							<rollingfile type="size" filename="./logs/tx_sync.log" maxsize="200000000" maxrolls="1"/>
						</outputs>
						<formats>
							<format id="default" format="%Date(2006-01-02 15:04:05.000)[%File:%Line]%Msg%n" />
						</formats>
					</seelog>`
)
