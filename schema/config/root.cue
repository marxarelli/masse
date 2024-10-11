package config

import (
	"wikimedia.org/dduvall/masse/common"
	"wikimedia.org/dduvall/masse/state"
	"wikimedia.org/dduvall/masse/target"
)

#Root: {
	parameters?: common.#Env
	chains!:     state.#Chains
	targets!:    target.#Targets
}
