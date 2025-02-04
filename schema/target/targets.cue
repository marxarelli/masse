package target

import (
	"github.com/marxarelli/masse/common"
	"github.com/marxarelli/masse/state"
)

#Runtime: {
	user:       string
	env:        common.#Env
	entrypoint: [string, ...string]
	arguments:  [...string]
	directory:  string
	stopSignal: string
}

#TargetPlatform: string | common.#Platform

#TargetDefaults: {
	platforms: [#TargetPlatform, ...#TargetPlatform] | *["linux/amd64"]
	labels:    common.#Labels                        | *null
	runtime:   {
		user:       #Runtime.user       | *"root"
		env:        #Runtime.env        | *null
		entrypoint: #Runtime.entrypoint | *null
		arguments:  #Runtime.arguments  | *null
		directory:  #Runtime.directory  | *"/"
		stopSignal: #Runtime.stopSignal | *"SIGTERM"
	}
}

#Target: {
	#default: #TargetDefaults

	build!:    state.#ChainRef
	platforms: #default.platforms
	labels:    #default.labels
	runtime:   {
		user:       #default.runtime.user
		env:        #default.runtime.env
		entrypoint: #default.runtime.entrypoint
		arguments:  #default.runtime.arguments
		directory:  #default.runtime.directory
		stopSignal: #default.runtime.stopSignal
	}

	if platforms != _|_ {
		platformsValue: [
			for p in platforms {
				[
					if (p & string) != _|_ {common.#Platform & {name: p}},
					p,
				][0]
			}
		]
	}
}

#Targets: {
	$default = #default: #TargetDefaults

	[=~"^[a-zA-Z_][a-zA-Z0-9_]*$"]: #Target & { #default: $default }
}
