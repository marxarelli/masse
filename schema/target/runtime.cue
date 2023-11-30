package target

import (
	"wikimedia.org/dduvall/phyton/schema/common"
)

#Runtime: {
	user: *(common.#User & {uid: 0}) | common.#User
	env:  common.#Env
	entrypoint?: [string, ...string]
	arguments?: [...string]
	directory:  *"/" | string
	stopSignal: *"SIGTERM" | string
}
