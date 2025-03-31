package state

import (
	"list"
	"strconv"
	"strings"
)

#Shell: {
	sh!:        string
	arguments?: string | [string, ...string] | *[]
	options?:   #RunOption | [#RunOption, ...#RunOption] | *[]

	let $arguments = list.FlattenN([*arguments | []], 1)
	let $options = list.FlattenN([*options | []], 1)

	let cmd = strings.TrimSpace(
		sh + " " + strings.Join([ for arg in $arguments { strconv.Quote(arg) } ], " ")
	)

	ops: [
		{
			run: ["/bin/sh", "-c", cmd]
			options: list.Concat([
				[
					{ customName: "💻 " + cmd },
				],
				$options,
			])
		}
	]
}
