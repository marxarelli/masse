package target

import (
	"wikimedia.org/dduvall/phyton/schema/common"
)

#Runtime: {
	user: *"root" | string
	env:  common.#Env
	entrypoint?: [string, ...string]
	arguments?: [...string]
	directory:  *"/" | string
	stopSignal: *"SIGTERM" | string
}