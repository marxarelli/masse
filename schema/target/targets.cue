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
#TargetPlatforms: [#TargetPlatform, ...#TargetPlatform]

#TargetDefaults: {
	platforms: #TargetPlatforms | *["linux/amd64"]
	labels:    common.#Labels                        | *null
	runtime:   {
		user:       #Runtime.user       | *"root"
		env:        #Runtime.env        | *null
		entrypoint: #Runtime.entrypoint | *null
		arguments:  #Runtime.arguments  | *null
		directory:  #Runtime.directory  | *"/"
		stopSignal: #Runtime.stopSignal | *"SIGTERM"
	}

	attestations: {
		sbom: {
			generator:  #SBOM.generator     | *""
			parameters: #SBOM.parameters    | *null
			scan:       #SBOM.scan          | *null
		}
	}
}

#Target: {
	#default: #TargetDefaults

	build!:    state.#ChainRef
	platforms: #TargetPlatforms | *#default.platforms
	labels:    common.#Labels   | *#default.labels
	runtime:   {
		user:       #Runtime.user       | *#default.runtime.user
		env:        #Runtime.env        | *#default.runtime.env
		entrypoint: #Runtime.entrypoint | *#default.runtime.entrypoint
		arguments:  #Runtime.arguments  | *#default.runtime.arguments
		directory:  #Runtime.directory  | *#default.runtime.directory
		stopSignal: #Runtime.stopSignal | *#default.runtime.stopSignal
	}

	attestations: {
		sbom: {
			generator:  #SBOM.generator  | *#default.attestations.sbom.generator
			parameters: #SBOM.parameters | *#default.attestations.sbom.parameters
			scan:       #SBOM.scan       | *#default.attestations.sbom.scan
		}
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
