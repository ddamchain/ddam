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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xchain/go-chain/crypto"
	"io/ioutil"
	"net/http"
)

type RemoteChainOpImpl struct {
	host string
	port int
	base string
	aop  accountOp
	show bool
}

// InitRemoteChainOp connect node by ip and port
func InitRemoteChainOp(ip string, port int, show bool, op accountOp) *RemoteChainOpImpl {
	ca := &RemoteChainOpImpl{
		aop:  op,
		show: show,
	}
	ca.Connect(ip, port)
	return ca
}

// Connect connect node by ip and port
func (ca *RemoteChainOpImpl) Connect(ip string, port int) error {
	if ip == "" {
		return nil
	}
	ca.host = ip
	ca.port = port
	ca.base = fmt.Sprintf("http://%v:%v", ip, port)
	return nil
}

func (ca *RemoteChainOpImpl) request(method string, params ...interface{}) *Result {
	if ca.base == "" {
		return opError(ErrUnConnected)
	}

	param := RPCReqObj{
		Method:  "Gx_" + method,
		Params:  params[:],
		ID:      1,
		Jsonrpc: "2.0",
	}

	if ca.show {
		fmt.Println("Request:")
		bs, _ := json.MarshalIndent(param, "", "\t")
		fmt.Println(string(bs))
		fmt.Println("==================================================================================")
	}

	paramBytes, err := json.Marshal(param)
	if err != nil {
		return opError(err)
	}

	resp, err := http.Post(ca.base, "application/json", bytes.NewReader(paramBytes))
	if err != nil {
		return opError(err)
	}
	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	ret := &RPCResObj{}
	if err := json.Unmarshal(responseBytes, ret); err != nil {
		return opError(err)
	}
	if ret.Error != nil {
		return opError(fmt.Errorf(ret.Error.Message))
	}
	return ret.Result
}

func (ca *RemoteChainOpImpl) nonce(addr string) (uint64, error) {
	ret := ca.request("nonce", addr)
	if !ret.IsSuccess() {
		return 0, fmt.Errorf(ret.Message)
	}
	return uint64(ret.Data.(float64)), nil
}

// Endpoint returns current connected ip and port
func (ca *RemoteChainOpImpl) Endpoint() string {
	return fmt.Sprintf("%v:%v", ca.host, ca.port)
}

// SendRaw send transaction to connected node
func (ca *RemoteChainOpImpl) SendRaw(tx *txRawData) *Result {
	r := ca.aop.AccountInfo()
	if !r.IsSuccess() {
		return r
	}
	aci := r.Data.(*Account)
	privateKey := crypto.HexToPrivateKey(aci.Sk)
	pubkey := crypto.HexToPubKey(aci.Pk)
	if privateKey.GetPubKey().Hex() != pubkey.Hex() {
		return opError(fmt.Errorf("privatekey or pubkey error"))
	}
	source := pubkey.GetAddress()
	if source.Hex() != aci.Address {
		return opError(fmt.Errorf("address error"))
	}

	if tx.Nonce == 0 {
		nonce, err := ca.nonce(aci.Address)
		if err != nil {
			return opError(err)
		}
		tx.Nonce = nonce
	}

	tranx := txRawToTransaction(tx)
	tranx.Hash = tranx.GenHash()
	sign, err := privateKey.Sign(tranx.Hash.Bytes())
	if err != nil {
		return opError(err)
	}
	tranx.Sign = &sign
	tx.Sign = sign.Hex()

	jsonByte, err := json.Marshal(tx)
	if err != nil {
		return opError(err)
	}

	ca.aop.(*AccountManager).resetExpireTime(aci.Address)
	// Signature is required here
	return ca.request("tx", string(jsonByte))
}

// Balance query Balance by address
func (ca *RemoteChainOpImpl) Balance(addr string) *Result {
	return ca.request("balance", addr)
}

// Nonce query Balance by address
func (ca *RemoteChainOpImpl) Nonce(addr string) *Result {
	return ca.request("nonce", addr)
}

func (ca *RemoteChainOpImpl) BlockHeight() *Result {
	return ca.request("blockHeight")
}

func (ca *RemoteChainOpImpl) TxInfo(hash string) *Result {
	return ca.request("transDetail", hash)
}

func (ca *RemoteChainOpImpl) BlockByHash(hash string) *Result {
	return ca.request("getBlockByHash", hash)
}

func (ca *RemoteChainOpImpl) BlockByHeight(h uint64) *Result {
	return ca.request("getBlockByHeight", h)
}

func (ca *RemoteChainOpImpl) TxReceipt(hash string) *Result {
	return ca.request("txReceipt", hash)
}

func (ca *RemoteChainOpImpl) Stake(addr string) *Result {
	return ca.request("stake", addr)
}
