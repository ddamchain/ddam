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
	"fmt"
	"github.com/xchain/go-chain/auth"
	"github.com/xchain/go-chain/common"
	"github.com/xchain/go-chain/consensus"
	"github.com/xchain/go-chain/crypto"
	"github.com/xchain/go-chain/global"
	"github.com/xchain/go-chain/global/types"
	"github.com/xchain/go-chain/network"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/xchain/go-chain/core"
)

var (
	cfg = &minerConfig{
		rpcLevel:      rpcLevelGx, // rpcLevelDev,
		rpcAddr:       "127.0.0.1",
		rpcPort:       8101,
		super:         false,
		testMode:      true,
		natIP:         "",
		natPort:       0,
		keystore:      "keystore",
		enableMonitor: false,
		chainID:       1,
		password:      "123",
	}
)

var (
	addr1      = "0x3056143f7287f11ad7de82bc90bd74b12bec157ae3e98fa6592650b774ef9eb5"
	addr2      = "0x966e3d026c59ce0f786d9cba7678de1d2f89d808e867556df47ac41fc5195baa"
	UMIDMgr    *auth.IdentityManager
	address    common.Address
	privateKey = "0x045c8153e5a849eef465244c0f6f40a43feaaa6855495b62a400cc78f9a6d61c76c09c3aaef393aa54bd2adc5633426e9645dfc36723a75af485c5f5c9f2c94658562fcdfb24e943cf257e25b9575216c6647c4e75e264507d2d57b3c8bc00b361"
	source     = ""
)

func initContext4Test() error {
	//global.Init("../../test/xconf_test.ini")
	//UMIDMgr = auth.NewIdentityManager()

	sk := crypto.HexToPrivateKey(privateKey)
	pk := sk.GetPubKey()
	address = pk.GetAddress()
	source = address.Hex()
	miner := &types.Miner{
		Addr:       address,
		PublicKey:  &pk,
		PrivateKey: sk,
	}
	global.Init("../test/xconf_test.ini")
	global.Context().Current = miner
	UMIDMgr = auth.NewIdentityManager()

	netCfg := network.NetworkConfig{
		IsSuper:         false,
		TestMode:        true,
		NatAddr:         "",
		NatPort:         0,
		SeedAddr:        "127.0.0.1",
		NodeIDHex:       address.Hex(),
		ChainID:         0,
		ProtocolVersion: 0,
		SeedIDs:         nil,
		PK:              pk.Hex(),
		SK:              sk.Hex(),
	}
	err := network.Init(netCfg)
	if err != nil {
		return err
	}

	// Init core block chain
	err = core.InitCore(UMIDMgr)
	if err != nil {
		return err
	}

	// Init consensus
	err = consensus.InitEngine(core.BlockChainImpl, UMIDMgr, network.GetNetInstance())
	return err
}

func TestRPC(t *testing.T) {
	err := initContext4Test()
	if err != nil {
		t.Fatalf("failed to initContext4Test")
	}
	gtxc := NewDdamApp()
	gtxc.config = cfg
	//new ConfFileManager
	cmr := types.NewConfINIManager("tas.ini")

	aop, err := newAccountOp("keystore")
	account := aop.NewAccount("123")
	aop.store.Close()
	addr := account.Data.(string)

	cmr.SetString(confSection, "miner", addr)

	//err = gtxc.fullInit()
	//if err != nil {
	//	t.Error(err)
	//}

	defer resetDb("testkey")
	cmr.Del(confSection, "miner")
	senderAddr := common.HexToAddress("0x593758f65911996c5bb2d143134430e2607e26ad815df6589455fb568b94e8b2")
	nonce := core.BlockChainImpl.GetNonce(senderAddr)
	fmt.Println("nonce :", nonce)

	core.BlockChainImpl.LatestStateDB().AddBalance(senderAddr, types.NewBigInt(100000000000).Value())
	fmt.Println("余额======》》", core.BlockChainImpl.GetBalance(senderAddr))
	privateKey := crypto.HexToPrivateKey("0x04db8949be0b4c733dc4c21dffdb647973634005e4b5beae38e1d264696a141c6cc6a7579c5d80d74766917f7ffa0e04a797933252bf0c818f0758d84b4ed0ab6e3d4e00f3ca3477c3260a1a94af9ce70947522372d097937ccae830b09748ffa4")

	//umid := []byte{137, 28, 54, 95, 251, 240, 164, 8, 198, 127, 154, 101, 31, 158, 179, 130, 176, 111, 248, 152, 203, 34, 1, 112, 192, 205, 96, 214, 183, 97, 194, 55} //001
	//umid := []byte{26, 199, 89, 153, 87, 215, 195, 183, 254, 45, 153, 234, 48, 62, 142, 38, 139, 86, 78, 216, 92, 255, 114, 99, 22, 16, 144, 253, 34, 12, 94, 252}  //002
	//umid := []byte{33, 165, 16, 69, 212, 221, 65, 20, 181, 209, 11, 7, 54, 54, 100, 94, 205, 156, 82, 133, 144, 61, 119, 5, 229, 139, 70, 215, 250, 129, 241, 153}   //003
	//umid := []byte{33, 165, 16, 69, 212, 221, 65, 20, 181, 209, 11, 7, 54, 54, 100, 94, 205, 156, 82, 133, 144, 61, 119, 5, 229, 139, 70, 215, 250, 129, 241, 153} //mac

	//a6  i5
	umid := []byte{234, 46, 131, 211, 67, 49, 160, 225, 0, 108, 142, 196, 150, 33, 224, 163, 193, 127, 249, 182, 204, 185, 139, 195, 102, 151, 51, 191, 91, 185, 154, 135}
	//a5 i7
	//umid := []byte{75, 62, 151, 59, 228, 193, 189, 109, 54, 23, 97, 107, 198, 67, 110, 100, 54, 93, 171, 136, 213, 219, 138, 83, 72, 63, 62, 85, 252, 163, 44, 164}
	//99k
	//umid := []byte{234, 218, 112, 2, 246, 78, 198, 137, 214, 4, 59, 6, 161, 178, 7, 7, 201, 52, 183, 137, 101, 225, 194, 160, 220, 238, 207, 62, 247, 220, 41, 171}

	//测试20t
	//umid := []byte{37, 254, 208, 142, 35, 250, 218, 173, 41, 212, 57, 5, 62, 187, 126, 3, 48, 105, 165, 224, 177, 183, 153, 231, 134, 109, 226, 107, 88, 91, 101, 53}
	//10t
	//umid:=[]byte{45, 246, 86, 180, 97, 58, 143, 219, 130, 43, 158, 230, 225, 148, 157, 27, 147, 153, 38, 220, 145, 25, 47, 69, 233, 34, 247, 95, 52, 129, 204, 134}
	//5t
	//umid:=[]byte{182, 32, 35, 224, 77, 181, 11, 127, 199, 205, 29, 81, 63, 232, 169, 11, 82, 165, 6, 74, 87, 255, 253, 157, 162, 71, 3 ,31, 50, 185, 48, 121}

	address := "0xa6c1106469d59abaca375504f232738ef6cc7750f1e330cb1708cb275e8965e1"
	//address:="0xa5258bd80d9e87aa4174decbe4aabd1f6c6689938ae52ea756395818c084887d"

	tx := &txRawData{Target: address, Value: 0, Gas: 10000, Gasprice: 10000, TxType: 1, Nonce: 1, Data: umid}

	//
	tranx := txRawToTransaction(tx)
	tranx.Hash = tranx.GenHash()
	sign, _ := privateKey.Sign(tranx.Hash.Bytes())
	tranx.Sign = &sign
	tx.Sign = sign.Hex()
	txdata, err := json.Marshal(tx)
	if err != nil {
		t.Error(err)
	}
	if err = gtxc.startRPC(core.BlockChainImpl, core.TxPool); err != nil {
		t.Error(err)
	}

	tests := []struct {
		method string
		params interface{}
	}{
		{"Gx_tx", string(txdata)},
		//{"Explorer_explorerBlockDetail", 1},
		//{"Dev_getTopBlock", nil},
		//{"Gx_balance", "0x135f237ff3616687e1fa0eb608f7dd5b7acba7f487ae78f8cfdc777663f6a3e0"},
		//{"Gx_blockHeight", nil},
		//{"Gx_nonce", "0x135f237ff3616687e1fa0eb608f7dd5b7acba7f487ae78f8cfdc777663f6a3e0"},
		//{"Gx_getBlockByHash", "0x9b2a5754771788fd013316171cd272971c8074ac78e0e75694e1094e7f0b6ecf"},
		//{"Gx_getBlockByHeight", 0},
		//{"Gx_transDetail", "0x135f237ff3616687e1fa0eb608f7dd5b7acba7f487ae78f8cfdc777663f6a3e0"},
		//{"Gx_txReceipt", "0x135f237ff3616687e1fa0eb608f7dd5b7acba7f487ae78f8cfdc777663f6a3e0"},
	}
	for _, test := range tests {
		fmt.Print("SIGN ADN DATA==========>", test)

		res, err := rpcPost(cfg.rpcAddr, uint(cfg.rpcPort), test.method, test.params)
		if err != nil {
			t.Errorf("%s failed: %v", test.method, err)
			return
		}
		if res.Error != nil {
			t.Errorf("%s failed: %v", test.method, res.Error.Message)
			return
		}
		data, _ := json.Marshal(res.Result)
		fmt.Print("+++++++DATA==========>", data)
		log.Printf("%s response data: %s", test.method, data)
	}
}

func resetDb(dbPath string) error {
	//core.BlockChainImpl.(*core.blockChainImpl).Close()
	//taslog.Close()
	fmt.Println("---reset db---")
	dir, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	for _, d := range dir {
		if d.IsDir() && strings.HasPrefix(d.Name(), "d_") {
			fmt.Printf("deleting folder: %s \n", d.Name())
			err = os.RemoveAll(d.Name())
			if err != nil {
				return err
			}
		}
		if d.IsDir() && strings.Compare(dbPath, d.Name()) == 0 {
			os.RemoveAll(d.Name())
		}

		if d.IsDir() && strings.Compare("logs", d.Name()) == 0 {
			os.RemoveAll(d.Name())
		}
	}
	os.RemoveAll(cfg.keystore)
	return nil
}

func TestMarshalTxRawData(t *testing.T) {
	tx := &txRawData{
		Target:   "0x123",
		Value:    100000000,
		Gas:      1304,
		Gasprice: 2324,
	}
	json, err := json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(json))

}

func TestUnmarhsalTxRawData(t *testing.T) {
	s := `{"target":"0x123","value":23,"gas":99,"gasprice":2324,"tx_type":0,"nonce":0,"data":"","sign":"","extra_data":""}`
	tx := &txRawData{}

	err := json.Unmarshal([]byte(s), tx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestByteToAddress(t *testing.T) {
	umid := []byte{234, 218, 112, 2, 246, 78, 198, 137, 214, 4, 59, 6, 161, 178, 7, 7, 201, 52, 183, 137, 101, 225, 194, 160, 220, 238, 207, 62, 247, 220, 41, 171} //99k
	//20t
	//umid := []byte{37, 254, 208, 142, 35, 250, 218, 173, 41, 212, 57, 5, 62, 187, 126, 3, 48, 105, 165, 224, 177, 183, 153, 231, 134, 109, 226, 107, 88, 91, 101, 53}
	a := common.BytesToHash(umid)

	fmt.Println(a.Hex())

	//k99k := "0xeada7002f64ec689d6043b06a1b20707c934b78965e1c2a0dceecf3ef7dc29ab"
	k20t := "0x25fed08e23fadaad29d439053ebb7e033069a5e0b1b799e7866de26b585b6535"
	fmt.Println(common.HexToHash(k20t).Bytes())
}
