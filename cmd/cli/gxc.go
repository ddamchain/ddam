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
	"os"
	"time"

	"github.com/xchain/go-chain/auth"
	"github.com/xchain/go-chain/consensus"
	"github.com/xchain/go-chain/crypto"
	"github.com/xchain/go-chain/global"
	"github.com/xchain/go-chain/global/types"
	"github.com/xchain/go-chain/xlog"

	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/core"
	"github.com/xchain/go-chain/network"
	"gopkg.in/alecthomas/kingpin.v2"

	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"strconv"
)

const (
	// confSection is default section configuration
	confSection     = "ddam"
	cliVersion      = 1
	protocolVersion = 1
)

type ddamApp struct {
	inited       bool
	account      Account
	config       *minerConfig
	rpcInstances []rpcApi
}

func (ddam *ddamApp) checkMiner(miner *types.Miner) error {
	if miner.PrivateKey.GetPubKey().Hex() != miner.PublicKey.Hex() {
		return fmt.Errorf("pk error")
	}
	if !miner.PublicKey.GetAddress().EqualTo(miner.Addr) {
		return fmt.Errorf("address error")
	}
	return nil
}

func (ddam *ddamApp) fullInit() error {
	var err error
	cfg := ddam.config

	// 1. init global context
	global.Init(cfg.confFile)

	// 2. get account
	conf := global.Context().Config.GetSectionManager(confSection)
	err = ddam.checkAddress(cfg.keystore, conf.GetString("miner", ""), cfg.password)
	if err != nil {
		return err
	}

	conf.SetString("miner", ddam.account.Address)
	showMsg("Your Miner Address:%s", ddam.account.Address)

	//set the time for proposer package
	timeForPackage := conf.GetInt("time_for_package", 2000)
	if timeForPackage > 100 && timeForPackage < 2000 {
		core.ProposerPackageTime = time.Duration(timeForPackage) * time.Millisecond
		showMsg("proposer uses the package config: timeForPackage %d ", timeForPackage)
	}

	//set the block gas limit for proposer package
	gasLimitForPackage := conf.GetInt("gas_limit_for_package", core.GasLimitPerBlock)
	if gasLimitForPackage > 10000 && gasLimitForPackage < core.GasLimitPerBlock {
		core.GasLimitForPackage = uint64(gasLimitForPackage)
		showMsg("proposer uses the package config: gasLimitForPackage %d ", gasLimitForPackage)
	}

	// Set current miner
	miner := &types.Miner{
		Addr:       common.HexToAddress(ddam.account.Address),
		PublicKey:  crypto.HexToPubKey(ddam.account.Pk),
		PrivateKey: crypto.HexToPrivateKey(ddam.account.Sk),
	}
	if err := ddam.checkMiner(miner); err != nil {
		return err
	}
	global.Context().Register("Current", miner)

	// Init network
	netCfg := network.NetworkConfig{
		IsSuper:         cfg.super,
		TestMode:        cfg.testMode,
		NatAddr:         cfg.natIP,
		NatPort:         cfg.natPort,
		SeedAddr:        cfg.seedIP,
		NodeIDHex:       miner.Addr.Hex(),
		ChainID:         cfg.chainID,
		ProtocolVersion: protocolVersion,
		SeedIDs:         []string{cfg.seedID},
		PK:              ddam.account.Pk,
		SK:              ddam.account.Sk,
	}
	err = network.Init(netCfg)
	if err != nil {
		return err
	}

	// Create umid manager
	umidMgr := auth.NewIdentityManager()

	// Init core block chain
	err = core.InitCore(umidMgr)
	if err != nil {
		return err
	}

	// Init consensus
	err = consensus.InitEngine(core.BlockChainImpl, umidMgr, network.GetNetInstance())
	if err != nil {
		return err
	}

	return nil
}

// miner start miner node
func (ddam *ddamApp) miner(cfg *minerConfig) error {
	ddam.config = cfg
	ddam.runtimeInit()
	err := ddam.fullInit()
	if err != nil {
		return err
	}
	if cfg.rpcEnable() {
		err = ddam.startRPC(core.BlockChainImpl, core.TxPool)
		if err != nil {
			return err
		}
	}
	showMsg("Syncing block info from ddam net, Waiting...")

	initMsgShower(core.BlockChainImpl)

	ddam.inited = true
	return nil
}

func (ddam *ddamApp) runtimeInit() {
	debug.SetGCPercent(100)
	debug.SetMaxStack(2 * 1000000000)
}

func (ddam *ddamApp) exit(ctrlC <-chan bool, quit chan<- bool) {
	<-ctrlC
	if core.BlockChainImpl == nil {
		return
	}
	showMsg("exiting...")
	core.BlockChainImpl.Close()
	xlog.Close()
	if ddam.inited {
		quit <- true
	} else {
		os.Exit(0)
	}
}

func (ddam *ddamApp) Run() {
	var err error

	// Control+c interrupt signal
	ctrlC := signals()
	quitChan := make(chan bool)
	go ddam.exit(ctrlC, quitChan)

	app := kingpin.New("Ddam", "A blockchain application.")
	app.HelpFlag.Short('h')
	configFile := app.Flag("config", "Config file").Default("xconf.ini").String()
	pprofPort := app.Flag("pprof", "enable pprof").Default("23333").Uint()
	keystore := app.Flag("keystore", "the keystore path, default is current path").Default("keystore").Short('k').String()

	// Console
	consoleCmd := app.Command("console", "start ddam console")
	showRequest := consoleCmd.Flag("show", "show the request json").Short('v').Bool()
	remoteHost := consoleCmd.Flag("host", "the node host address to connect").Short('i').String()
	remotePort := consoleCmd.Flag("port", "the node host port to connect").Short('p').Default("8101").Int()

	// Version
	versionCmd := app.Command("version", "show ddam cliVersion")

	// Mine
	mineCmd := app.Command("miner", "miner start")

	// Rpc analysis
	rpc := mineCmd.Flag("rpc", "start rpc server and specify the rpc service level").Default(strconv.FormatInt(int64(rpcLevelNone), 10)).Int()
	enableMonitor := mineCmd.Flag("monitor", "enable monitor").Default("false").Bool()
	addrRPC := mineCmd.Flag("rpcaddr", "rpc service host").Short('r').Default("0.0.0.0").IP()
	rpcServicePort := mineCmd.Flag("rpcport", "rpc service port").Short('p').Default("8101").Uint16()
	super := mineCmd.Flag("super", "start super node").Bool()
	passWd := mineCmd.Flag("password", "login password").Default("123").String()

	// In test mode, P2P NAT is closed
	testMode := mineCmd.Flag("test", "test mode").Bool()
	seedAddr := mineCmd.Flag("seed", "seed address").String()
	seedID := mineCmd.Flag("seedid", "seed ID").String()
	natAddr := mineCmd.Flag("nat", "nat server IP").String()
	natPort := mineCmd.Flag("natport", "nat server port").Default("0").Uint16()
	chainID := mineCmd.Flag("chainid", "chain ID").Default("0").Uint16()

	command, err := app.Parse(os.Args[1:])
	if err != nil {
		kingpin.Fatalf("%s, try --help", err)
	}

	switch command {
	case versionCmd.FullCommand():
		showMsg("Ddam Version:%d", cliVersion)
		os.Exit(0)
	case consoleCmd.FullCommand():
		err := ConsoleInit(*keystore, *remoteHost, *remotePort, *showRequest)
		if err != nil {
			showMsg(err.Error())
			os.Exit(0)
		}
	case mineCmd.FullCommand():
		go func() {
			http.ListenAndServe(fmt.Sprintf(":%d", *pprofPort), nil)
			runtime.SetBlockProfileRate(1)
			runtime.SetMutexProfileFraction(1)
		}()

		cfg := &minerConfig{
			confFile:      *configFile,
			rpcLevel:      rpcLevel(*rpc),
			rpcAddr:       addrRPC.String(),
			rpcPort:       *rpcServicePort,
			super:         *super,
			testMode:      *testMode,
			natIP:         *natAddr,
			natPort:       *natPort,
			seedIP:        *seedAddr,
			seedID:        *seedID,
			keystore:      *keystore,
			enableMonitor: *enableMonitor,
			chainID:       *chainID,
			password:      *passWd,
		}

		// Start miner
		err := ddam.miner(cfg)
		if err != nil {
			showMsg("miner start error:%v", err)
			os.Exit(0)
		}
	}
	<-quitChan
}

func (ddam *ddamApp) checkAddress(keystore, address, password string) error {
	aop, err := initAccountManager(keystore, true)
	if err != nil {
		return err
	}
	defer aop.Close()

	acm := aop.(*AccountManager)
	if address != "" {
		aci, err := acm.checkMinerAccount(address, password)
		if err != nil {
			return fmt.Errorf("cannot get miner, err:%v", err.Error())
		}
		ddam.account = aci.Account
		return nil
	}
	acc, _ := acm.getFirstMinerAccount(password)
	if acc != nil {
		ddam.account = *acc
		return nil
	}

	return fmt.Errorf("please provide a miner account and correct password! ")
}

func NewDdamApp() *ddamApp {
	return &ddamApp{}
}
