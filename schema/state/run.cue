package state

import (
	"list"
	"strings"
)

#Run: {
	run!:       string
	arguments?: string | [string, ...string]
	options?:   #RunOption | [#RunOption, ...#RunOption]

	#defaultOptions: [
		{ customName: strings.TrimSpace("💻 " + run + " " + strings.Join(list.FlattenN([*arguments | []], 1), " ")) },
	]
}
