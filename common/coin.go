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
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
**  Creator: pxf
**  Date: 2019/1/8 下午3:33
**  Description:
 */

const (
	AM   uint64 = 1
	KAM         = 1000
	MAM         = 1000000
	DDAM        = 1000000000
)

var (
	ErrEmptyStr   = fmt.Errorf("empty string")
	ErrIllegalStr = fmt.Errorf("illegal gasprice string")
)

var am, _ = regexp.Compile("^([0-9]+)(am|kam|mam|ddam)$")

// ParseCoin parses string to amount
func ParseCoin(s string) (uint64, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return 0, ErrEmptyStr
	}

	arr := am.FindAllStringSubmatch(s, -1)
	if arr == nil || len(arr) == 0 {
		return 0, ErrIllegalStr
	}
	ret := arr[0]
	if ret == nil || len(ret) != 3 {
		return 0, ErrIllegalStr
	}
	num, err := strconv.Atoi(ret[1])
	if err != nil {
		return 0, err
	}
	unit := AM
	if len(ret) == 3 {
		switch ret[2] {
		case "kam":
			unit = KAM
		case "mam":
			unit = MAM
		case "tas":
			unit = DDAM
		}
	}
	//fmt.Println(am.FindAllString(s, -1))
	return uint64(num) * unit, nil
}

func DDAM2AM(v uint64) uint64 {
	return v * DDAM
}

func Value2AM(v float64) uint64 {
	return uint64(v * float64(DDAM))
}

func AM2DDAM(v uint64) float64 {
	return float64(v) / float64(DDAM)
}
