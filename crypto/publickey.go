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

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/xchain/go-chain/common"
	"io"

	"github.com/xchain/go-chain/crypto/ecies"
	"github.com/xchain/go-chain/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
)

// PublicKey data struct
type PublicKey struct {
	PubKey ecdsa.PublicKey
}

// Verify the validation of signature and message
func (pk PublicKey) Verify(hash []byte, s *Sign) bool {
	return secp256k1.VerifySignature(pk.Bytes(), hash, s.Bytes()[:64])
}

// GetAddress obtains the address mapped from the public key
func (pk PublicKey) GetAddress() common.Address {
	x := pk.PubKey.X.Bytes()
	y := pk.PubKey.Y.Bytes()
	x = append(x, y...)

	addrBuf := sha3.Sum256(x)
	return common.BytesToAddress(addrBuf[:])
}

// Bytes converts the public key to a byte array
func (pk PublicKey) Bytes() []byte {
	buf := elliptic.Marshal(pk.PubKey.Curve, pk.PubKey.X, pk.PubKey.Y)
	//fmt.Printf("end pub key marshal, len=%v, data=%v\n", len(buf), buf)
	return buf
}

// BytesToPublicKey returns a public key with the byte array imported
func BytesToPublicKey(data []byte) (pk *PublicKey) {
	pk = new(PublicKey)
	pk.PubKey.Curve = getDefaultCurve()
	//fmt.Printf("begin pub key unmarshal, len=%v, data=%v.\n", len(data), data)
	x, y := elliptic.Unmarshal(pk.PubKey.Curve, data)
	if x == nil || y == nil {
		panic("unmarshal public key failed.")
	}
	pk.PubKey.X = x
	pk.PubKey.Y = y
	return
}

// Hex converts the public key to a hex string
func (pk PublicKey) Hex() string {
	return common.ToHex(pk.Bytes())
}

// Encrypt returns the cipher text of the message
func (pk *PublicKey) Encrypt(rand io.Reader, msg []byte) ([]byte, error) {
	return Encrypt(rand, pk, msg)
}

// HexToPubKey returns a public key with the hex string imported
func HexToPubKey(s string) (pk *PublicKey) {
	if len(s) < len(common.HexPrefix) || s[:len(common.HexPrefix)] != common.HexPrefix {
		return
	}
	pk = BytesToPublicKey(common.FromHex(s))
	return
}

// Encrypt returns the cipher text of the message
func Encrypt(rand io.Reader, pub *PublicKey, msg []byte) (ct []byte, err error) {
	pubECIES := ecies.ImportECDSAPublic(&pub.PubKey)
	return ecies.Encrypt(rand, pubECIES, msg, nil, nil)
}
