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

package common

import (
	"bytes"
	"crypto/sha256"
	"hash"
	"sync"

	"github.com/moonfruit/go-shabal"
)

var sha256Pool = sync.Pool{
	New: func() interface{} {
		return sha256.New()
	},
}

var shabalPool = sync.Pool{
	New: func() interface{} {
		return shabal.NewShabal256()
	},
}

// Sha256 computes Sha256 value of the input byte array
func Sha256(blockByte []byte) []byte {
	hasher := sha256Pool.Get().(hash.Hash)
	hasher.Reset()
	defer sha256Pool.Put(hasher)

	hasher.Write(blockByte)
	return hasher.Sum(nil)

}

// BytesCombine combines multiple byte arrays into one byte array
func BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}

// Shabal256 computes Shabal256 value of the input byte array
func Shabal256(bytes []byte) []byte {
	hasher := shabalPool.Get().(hash.Hash)
	hasher.Reset()
	defer shabalPool.Put(hasher)

	hasher.Write(bytes)
	return hasher.Sum(nil)

}
