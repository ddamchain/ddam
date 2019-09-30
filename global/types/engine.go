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

// Package auth provides functions on authority management
package types

type AuthCode []byte

type UMID []byte

const UMIDLength = 32

func (ac AuthCode) Bytes() []byte {
	return ac
}

func BytesToAuthCode(bs []byte) AuthCode {
	return bs
}

//func BytesToUMID(b []byte) UMID {
//	var a UMID
//	a.SetBytes(b)
//	return a
//}

//func (a *UMID) SetBytes(b []byte) {
//	if len(b) > len(a) {
//		b = b[len(b)-UMIDLength:]
//	}
//	copy(a[:], b[:])
//}

// ConsensusEngine are consensus block verifier
type ConsensusEngine interface {
	// check the new block
	// mainly verify the cast legality
	VerifyBlockHeader(bh *BlockHeader) (bool, error)

	VerifyBlockHeaderPair(prevBH, bh *BlockHeader) (bool, error)

	AvgBlockTime() int
}
