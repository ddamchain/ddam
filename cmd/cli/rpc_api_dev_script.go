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
	"github.com/xchain/go-chain/crypto"

	"github.com/xchain/go-chain/common"
)

func (api *RpcDevImpl) ScriptTransferTx(privateKey string, from string, to string, amount uint64, nonce uint64, txType int, gasPrice uint64) (*Result, error) {
	return api.TxUnSafe(privateKey, to, amount, gasPrice, gasPrice, nonce, txType, "")
}

// TxUnSafe sends a transaction by submitting the privateKey.
// It is not safe for users, used for testing purpose
func (api *RpcDevImpl) TxUnSafe(privateKey, target string, value, gas, gasprice, nonce uint64, txType int, data string) (*Result, error) {
	txRaw := &txRawData{
		Target:   target,
		Value:    common.DDAM2AM(value),
		Gas:      gas,
		Gasprice: gasprice,
		Nonce:    nonce,
		TxType:   txType,
		Data:     []byte(data),
	}
	sk := crypto.HexToPrivateKey(privateKey)
	if sk == nil {
		return failResult(fmt.Sprintf("parse private key fail:%v", privateKey))
	}
	trans := txRawToTransaction(txRaw)
	trans.Hash = trans.GenHash()
	sign, err := sk.Sign(trans.Hash.Bytes())
	if err != nil {
		failResult(err.Error())
	}
	trans.Sign = &sign

	if err := api.sendTransaction(trans); err != nil {
		return failResult(err.Error())
	}
	return successResult(trans.Hash.Hex())
}
