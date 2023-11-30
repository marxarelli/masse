package config

import (
	"wikimedia.org/dduvall/phyton/schema/common"
	"wikimedia.org/dduvall/phyton/schema/state"
	"wikimedia.org/dduvall/phyton/schema/target"
)

#Root: {
	parameters?: common.#Env
	chains!:     state.#Chains
	targets!:    target.#Targets
}
