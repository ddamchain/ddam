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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/peterh/liner"
	"github.com/xchain/go-chain/common"
)

type baseCmd struct {
	name string
	help string
	fs   *flag.FlagSet
}

func genBaseCmd(n string, h string) *baseCmd {
	return &baseCmd{
		name: n,
		help: h,
		fs:   flag.NewFlagSet(n, flag.ContinueOnError),
	}
}

type newAccountCmd struct {
	baseCmd
	password string
	miner    bool
}

func output(msg ...interface{}) {
	fmt.Println(msg...)
}

func genNewAccountCmd() *newAccountCmd {
	c := &newAccountCmd{
		baseCmd: *genBaseCmd("newaccount", "create account"),
	}
	c.fs.StringVar(&c.password, "password", "", "password for the account")
	return c
}

func (c *newAccountCmd) parse(args []string) bool {
	err := c.fs.Parse(args)
	if err != nil {
		output(err.Error())
		return false
	}
	pass := strings.TrimSpace(c.password)
	if len(pass) == 0 {
		output("Please input password")
		return false
	}
	if len(pass) > 50 || len(pass) < 3 {
		output("password length should between 3-50")
		return false
	}
	return true
}

type unlockCmd struct {
	baseCmd
	addr     string
	duration uint
}

func genUnlockCmd() *unlockCmd {
	c := &unlockCmd{
		baseCmd: *genBaseCmd("unlock", "unlock the account"),
	}
	c.fs.StringVar(&c.addr, "addr", "", "the account address")
	c.fs.UintVar(&c.duration, "duration", 120, "unlock duration, default 120 secs")
	return c
}

func (c *unlockCmd) parse(args []string) bool {
	if err := c.fs.Parse(args); err != nil {
		output(err.Error())
		return false
	}
	if strings.TrimSpace(c.addr) == "" {
		output("please input the address")
		c.fs.PrintDefaults()
		return false
	}

	if !validateAddress(c.addr) {
		output("Wrong address format")
		return false
	}
	return true
}

type balanceCmd struct {
	baseCmd
	addr string
}

func genBalanceCmd() *balanceCmd {
	c := &balanceCmd{
		baseCmd: *genBaseCmd("balance", "get the balance of the current unlocked account"),
	}
	c.fs.StringVar(&c.addr, "addr", "", "the account address")
	return c
}

func (c *balanceCmd) parse(args []string) bool {
	if err := c.fs.Parse(args); err != nil {
		output(err.Error())
		return false
	}
	if strings.TrimSpace(c.addr) == "" {
		output("please input the address")
		c.fs.PrintDefaults()
		return false
	}
	if !validateAddress(c.addr) {
		output("Wrong address format")
		return false
	}
	return true
}

type stakeCmd struct {
	baseCmd
	addr string
}

func genStakeCmd() *stakeCmd {
	c := &stakeCmd{
		baseCmd: *genBaseCmd("stake", "get the stake of the current unlocked account"),
	}
	c.fs.StringVar(&c.addr, "addr", "", "the account address")
	return c
}

func (c *stakeCmd) parse(args []string) bool {
	if err := c.fs.Parse(args); err != nil {
		output(err.Error())
		return false
	}
	if strings.TrimSpace(c.addr) == "" {
		output("please input the address")
		c.fs.PrintDefaults()
		return false
	}
	if !validateAddress(c.addr) {
		output("Wrong address format")
		return false
	}
	return true
}

type nonceCmd struct {
	baseCmd
	addr string
}

func genNonceCmd() *nonceCmd {
	c := &nonceCmd{
		baseCmd: *genBaseCmd("nonce", "get the nonce of the current unlocked account"),
	}
	c.fs.StringVar(&c.addr, "addr", "", "the account address")
	return c
}

func (c *nonceCmd) parse(args []string) bool {
	if err := c.fs.Parse(args); err != nil {
		output(err.Error())
		return false
	}
	if strings.TrimSpace(c.addr) == "" {
		output("please input the address")
		c.fs.PrintDefaults()
		return false
	}
	if !validateAddress(c.addr) {
		output("Wrong address format")
		return false
	}
	return true
}

type connectCmd struct {
	baseCmd
	host string
	port int
}

func genConnectCmd() *connectCmd {
	c := &connectCmd{
		baseCmd: *genBaseCmd("connect", "connect to one tas node"),
	}
	c.fs.StringVar(&c.host, "host", "", "the node ip")
	c.fs.IntVar(&c.port, "port", 8101, "the node port, default is 8101")
	return c
}

func (c *connectCmd) parse(args []string) bool {
	if err := c.fs.Parse(args); err != nil {
		output(err.Error())
		return false
	}
	if strings.TrimSpace(c.host) == "" {
		output("please input the host,available testnet hosts are node1.taschain.cn,node2.taschain.cn,node3.taschain.cn,node4.taschain.cn,node5.taschain.cn")
		c.fs.PrintDefaults()
		return false
	}
	if c.port == 0 {
		output("please input the port")
		c.fs.PrintDefaults()
		return false
	}
	return true
}

type txCmd struct {
	baseCmd
	hash     string
	executed bool
}

func genTxCmd() *txCmd {
	c := &txCmd{
		baseCmd: *genBaseCmd("tx", "get transaction detail"),
	}
	c.fs.StringVar(&c.hash, "hash", "", "the hex transaction hash")
	c.fs.BoolVar(&c.executed, "executed", false, "get executed transaction detail")
	return c
}

func (c *txCmd) parse(args []string) bool {
	if err := c.fs.Parse(args); err != nil {
		output(err.Error())
		return false
	}
	if strings.TrimSpace(c.hash) == "" {
		output("please input the transaction hash")
		c.fs.PrintDefaults()
		return false
	}
	if !validateHash(c.hash) {
		output("Wrong hash format")
		return false
	}
	return true
}

type blockCmd struct {
	baseCmd
	hash   string
	height uint64
}

func genBlockCmd() *blockCmd {
	c := &blockCmd{
		baseCmd: *genBaseCmd("block", "get block detail"),
	}
	c.fs.StringVar(&c.hash, "hash", "", "the hex block hash")
	c.fs.Uint64Var(&c.height, "height", 0, "the block height")
	return c
}

func (c *blockCmd) parse(args []string) bool {
	if err := c.fs.Parse(args); err != nil {
		output(err.Error())
		return false
	}
	if len(c.hash) > 0 {
		if !validateHash(c.hash) {
			output("Wrong hash format")
			return false
		}
	}
	return true
}

type gasBaseCmd struct {
	baseCmd
	gaslimit    uint64
	gasPriceStr string
	gasPrice    uint64
}

func genGasBaseCmd(n string, h string) *gasBaseCmd {
	c := &gasBaseCmd{
		baseCmd: *genBaseCmd(n, h),
	}
	return c
}

func (c *gasBaseCmd) parseGasPrice() bool {
	gp, err := common.ParseCoin(c.gasPriceStr)
	if err != nil {
		output(fmt.Sprintf("%v:%v, correct example: 100AM,100kAM,1mAM,1DDAM", err, c.gasPriceStr))
		return false
	}
	c.gasPrice = gp
	return true
}

func (c *gasBaseCmd) initBase() {
	c.fs.Uint64Var(&c.gaslimit, "gaslimit", 3000, "gas limit, default 3000")
	c.fs.StringVar(&c.gasPriceStr, "gasprice", "500AM", "gas price, default 500AM")
}

type sendTxCmd struct {
	gasBaseCmd
	to     string
	value  float64
	data   string
	nonce  uint64
	txType int
}

func genSendTxCmd() *sendTxCmd {
	c := &sendTxCmd{
		gasBaseCmd: *genGasBaseCmd("sendtx", "send a transaction to the tas system"),
	}
	c.initBase()
	c.fs.StringVar(&c.to, "to", "", "the transaction receiver address")
	c.fs.Float64Var(&c.value, "value", 0.0, "transfer value in tas unit")
	c.fs.StringVar(&c.data, "data", "", "transaction data")
	c.fs.Uint64Var(&c.nonce, "nonce", 0, "nonce, optional. will use default nonce on core if not specified")
	c.fs.IntVar(&c.txType, "type", 0, "transaction type: 0=general tx, 1=bind, 2=transfer bind, 3=unbind, 4=stake add, 5=stake reduce")
	return c
}

func (c *sendTxCmd) toTxRaw() *txRawData {
	return &txRawData{
		Target:   c.to,
		Value:    common.Value2AM(c.value),
		TxType:   c.txType,
		Data:     []byte(c.data),
		Gas:      c.gaslimit,
		Gasprice: c.gasPrice,
		Nonce:    c.nonce,
	}
}

func (c *sendTxCmd) parse(args []string) bool {
	if err := c.fs.Parse(args); err != nil {
		output(err.Error())
		return false
	}
	if !validateValue(c.value) {
		output("please input the value greater than 0")
		return false
	}
	if !validateTxType(c.txType) {
		output("Not supported transaction type")
		return false
	}
	if strings.TrimSpace(c.to) == "" {
		output("please input the target address")
		c.fs.PrintDefaults()
		return false
	} else {
		if !validateAddress(strings.TrimSpace(c.to)) {
			output("Wrong address format")
			return false
		}
	}

	if !c.parseGasPrice() {
		return false
	}

	return true
}

type importKeyCmd struct {
	baseCmd
	key      string
	password string
	miner    bool
}

func genImportKeyCmd() *importKeyCmd {
	c := &importKeyCmd{
		baseCmd: *genBaseCmd("importkey", "import private key"),
	}
	c.fs.StringVar(&c.key, "privatekey", "", "private key imported for the account")
	c.fs.StringVar(&c.password, "password", "", "password for the account")
	c.fs.BoolVar(&c.miner, "miner", false, "create the account for miner if set")
	return c
}

func (c *importKeyCmd) parse(args []string) bool {
	err := c.fs.Parse(args)
	if err != nil {
		output(err.Error())
		return false
	}
	key := strings.TrimSpace(c.key)
	if len(key) == 0 {
		output("Please input private key")
		return false
	}
	if !validateKey(key) {
		output("Private key is invalid")
		return false
	}
	if !validateHash(key) {
		output("Private key is invalid")
		return false
	}
	pass := strings.TrimSpace(c.password)
	if len(pass) == 0 {
		output("Please input password")
		return false
	}
	if len(pass) > 50 || len(pass) < 3 {
		output("password length should between 3-50")
		return false
	}
	return true
}

var cmdNewAccount = genNewAccountCmd()
var cmdExit = genBaseCmd("exit", "quit  ddam")
var cmdHelp = genBaseCmd("help", "show help info")
var cmdAccountList = genBaseCmd("accountlist", "list the account of the keystore")
var cmdUnlock = genUnlockCmd()
var cmdBalance = genBalanceCmd()
var cmdNonce = genNonceCmd()
var cmdAccountInfo = genBaseCmd("accountinfo", "get the info of the current unlocked account")
var cmdDelAccount = genBaseCmd("delaccount", "delete the info of the current unlocked account")
var cmdConnect = genConnectCmd()
var cmdBlockHeight = genBaseCmd("blockheight", "the current block height")
var cmdStake = genStakeCmd()
var cmdTx = genTxCmd()
var cmdBlock = genBlockCmd()
var cmdSendTx = genSendTxCmd()

var cmdImportKey = genImportKeyCmd()

//var cmdExportKey = genExportKeyCmd()
var cmdExportKey = genBaseCmd("exportkey", "export private key")

var list = make([]*baseCmd, 0)

func init() {
	list = append(list, cmdHelp)
	list = append(list, &cmdNewAccount.baseCmd)
	list = append(list, cmdAccountList)
	list = append(list, &cmdUnlock.baseCmd)
	list = append(list, &cmdBalance.baseCmd)
	list = append(list, &cmdNonce.baseCmd)
	list = append(list, cmdAccountInfo)
	list = append(list, &cmdStake.baseCmd)
	list = append(list, cmdDelAccount)
	list = append(list, &cmdConnect.baseCmd)
	list = append(list, cmdBlockHeight)
	list = append(list, &cmdTx.baseCmd)
	list = append(list, &cmdBlock.baseCmd)
	list = append(list, &cmdSendTx.baseCmd)
	list = append(list, &cmdImportKey.baseCmd)
	list = append(list, cmdExportKey)
	list = append(list, cmdExit)
}

func Usage() {
	output("Usage:")
	for _, cmd := range list {
		output(" " + cmd.name + ":\t" + cmd.help)
		cmd.fs.PrintDefaults()
		fmt.Print("\n")
	}
}

func ConsoleInit(keystore, host string, port int, show bool) error {
	aop, err := initAccountManager(keystore, false)
	if err != nil {
		return err
	}
	chainop := InitRemoteChainOp(host, port, show, aop)
	if chainop.base != "" {

	}

	loop(aop, chainop)

	return nil
}

func handleCmd(handle func() *Result) {
	ret := handle()
	if !ret.IsSuccess() {
		output(ret.Message)
	} else {
		bs, err := json.MarshalIndent(ret, "", "\t")
		if err != nil {
			output(err.Error())
		} else {
			output(string(bs))
		}
	}
}

func unlockLoop(cmd *unlockCmd, acm accountOp) {
	c := 0

	for c < 3 {
		c++

		bs, err := gopass.GetPasswdPrompt("please input password: ", true, os.Stdin, os.Stdout)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		ret := acm.UnLock(cmd.addr, string(bs), cmd.duration)
		if ret.IsSuccess() {
			fmt.Printf("unlock will last %v secs:%v\n", cmd.duration, cmd.addr)
			break
		} else {
			fmt.Fprintln(os.Stderr, ret.Message)
		}
	}
}

func loop(acm accountOp, chainOp chainOp) {

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	items := make([]string, len(list))
	for idx, cmd := range list {
		items[idx] = cmd.name
	}

	line.SetCompleter(func(line string) (c []string) {
		for _, n := range items {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	for {
		ep := chainOp.Endpoint()
		if ep == ":0" {
			ep = "not connected"
		}
		input, err := line.Prompt(fmt.Sprintf("ddam:%v > ", ep))
		if err != nil {
			if err == liner.ErrPromptAborted {
				line.Close()
				os.Exit(0)
			}
			fmt.Fprintln(os.Stderr, err)
		}

		inputArr, err := parseCommandLine(input)
		if err != nil {
			fmt.Printf("%s", err.Error())
		}

		line.AppendHistory(input)

		if len(inputArr) == 0 {
			continue
		}
		cmdStr := inputArr[0]
		args := inputArr[1:]

		switch cmdStr {
		case "":
			break
		case cmdNewAccount.name:
			cmd := genNewAccountCmd()
			if cmd.parse(args) {
				handleCmd(func() *Result {
					return acm.NewAccount(cmd.password)
				})
			}
		case cmdExit.name, "quit":
			fmt.Printf("thank you, bye\n")
			line.Close()
			os.Exit(0)
		case cmdHelp.name:
			Usage()
		case cmdAccountList.name:
			handleCmd(func() *Result {
				return acm.AccountList()
			})
		case cmdUnlock.name:
			cmd := genUnlockCmd()
			if cmd.parse(args) {
				unlockLoop(cmd, acm)
			}
		case cmdAccountInfo.name:
			handleCmd(func() *Result {
				return acm.AccountInfo()
			})
		case cmdDelAccount.name:
			handleCmd(func() *Result {
				return acm.DeleteAccount()
			})
		case cmdConnect.name:
			cmd := genConnectCmd()
			if cmd.parse(args) {
				chainOp.Connect(cmd.host, cmd.port)
			}

		case cmdBalance.name:
			cmd := genBalanceCmd()
			if cmd.parse(args) {
				handleCmd(func() *Result {
					return chainOp.Balance(cmd.addr)
				})
			}
		case cmdNonce.name:
			cmd := genNonceCmd()
			if cmd.parse(args) {
				handleCmd(func() *Result {
					return chainOp.Nonce(cmd.addr)
				})
			}
		case cmdBlockHeight.name:
			handleCmd(func() *Result {
				return chainOp.BlockHeight()
			})
		case cmdStake.name:
			cmd := genStakeCmd()
			if cmd.parse(args) {
				handleCmd(func() *Result {
					return chainOp.Stake(cmd.addr)
				})
			}
		case cmdTx.name:
			cmd := genTxCmd()
			if cmd.parse(args) {
				handleCmd(func() *Result {
					if cmd.executed {
						return chainOp.TxReceipt(cmd.hash)
					}
					return chainOp.TxInfo(cmd.hash)
				})
			}
		case cmdBlock.name:
			cmd := genBlockCmd()
			if cmd.parse(args) {
				handleCmd(func() *Result {
					if cmd.hash != "" {
						return chainOp.BlockByHash(cmd.hash)
					}
					return chainOp.BlockByHeight(cmd.height)
				})
			}
		case cmdSendTx.name:
			cmd := genSendTxCmd()
			if cmd.parse(args) {
				handleCmd(func() *Result {
					return chainOp.SendRaw(cmd.toTxRaw())
				})
			}
		case cmdImportKey.name:
			cmd := genImportKeyCmd()
			if cmd.parse(args) {
				handleCmd(func() *Result {
					return acm.NewAccountByImportKey(cmd.key, cmd.password)
				})
			}
		case cmdExportKey.name:
			handleCmd(func() *Result {
				return acm.ExportKey()
			})
		default:
			fmt.Printf("not supported command %v\n", cmdStr)
			Usage()
		}
	}
}

func parseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New(fmt.Sprintf("Unclosed quote in command line: %s", command))
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}
