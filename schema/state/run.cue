package state

import (
	"list"
)

#Run: {
	run!: string
	arguments?: [string, ...string]
	options?: #RunOption | [#RunOption, ...#RunOption]
	if options != _|_ {
		optionsValue: list.FlattenN([options], 1)
	}
}
