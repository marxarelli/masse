package layout

import (
	"wikimedia.org/dduvall/phyton/schema/common"
)

#ImageConfig: {
	from:        string
	user:        *(common.#User & {uid: 0}) | common.#User
	environment: common.#Env
	entrypoint?: [string, ...string]
	defaultArguments?: [string, ...string]
	workingDirectory: *"/" | string
	labels:           common.#Labels
	stopSignal:       string
}
