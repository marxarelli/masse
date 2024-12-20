package npm

import (
	"list"
	"github.com/marxarelli/masse/state"
)

install: {
	#command:     "ci" | *"install"
	#only:        string | *""
	#cache:       string | *"/var/lib/cache/npm"
	#options?:    state.#RunOption | [state.#RunOption, ...state.#RunOption]

	let $options = list.Concat([
		list.FlattenN(
			[
				if #options != _|_ {
					#options
				}
			],
			1,
		),
		[
			{ env: { NPM_CONFIG_CACHE: #cache } },
			{ cache: #cache, access: "locked" },
		],
	])

	let flags = [
		if #only != "" {
			" --only=\(#only)"
		},
		"",
	][0]

	{
		sh: "npm \(#command)\(flags) && npm dedupe"
		options: $options
	}
}
