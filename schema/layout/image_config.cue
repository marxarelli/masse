package layout

import (
	"wikimedia.org/dduvall/phyton/common"
)

#ImageConfig: {
	user: common.#User
	exposed_ports: {}
	environment: common.#Env
	entrypoint: [...string]
	default_arguments: [...string]
	working_directory: string
	labels: [...common.#Label]
	stop_signal: string
}
