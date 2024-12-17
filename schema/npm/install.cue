package npm

import (
	"list"
	"wikimedia.org/dduvall/masse/state"
)

install: {
	#command:     "ci" | *"install"
	#only:        string | *""
	#cache:       string | *"/var/lib/cache/npm"
	#options?:    state.#RunOption | [state.#RunOption, ...state.#RunOption]

	let $options = list.Concat([
		[
			{ env: { NPM_CONFIG_CACHE: #cache } },
			{ cache: #cache, access: "locked" },
		],
		list.FlattenN([#options | []], 1),
	])

	state.#Ops & {
		ops: [
			{
				run: "npm \(#command)"
				if #only != "" {
					arguments: ["--only=\(#only)"]
				}
				options: $options
			},
			{
				run: "npm dedupe"
				options: $options
			}
		]
	}
}
