package common

import (
	"list"
	"strings"
)

#Mode: {
	#NumericMode | #SymbolicMode
	value: uint32 & (>=0 & <=0o7777)
}

#NumericMode: {
	mode!: uint32
	value: mode
}

#SymbolicMode: {
	mode!: =~"^[r-][w-][xsS-][r-][w-][xsS-][r-][w-][xtT-]$"
	_map: [
		{"-": 0, r: 0o0400},
		{"-": 0, w: 0o0200},
		{"-": 0, x: 0o0100, s: 0o4100, S: 0o4000},
		{"-": 0, r: 0o0040},
		{"-": 0, w: 0o0020},
		{"-": 0, x: 0o0010, s: 0o2010, S: 0o2000},
		{"-": 0, r: 0o0004},
		{"-": 0, w: 0o0002},
		{"-": 0, x: 0o0001, t: 0o1001, T: 0o1000},
	]

	value: list.Sum([ for i, r in strings.Split(mode, "") {_map[i][r]}])
}
