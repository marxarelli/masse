package config

import (
	"wikimedia.org/releng/masse/common"
	"wikimedia.org/releng/masse/state"
	"wikimedia.org/releng/masse/target"
)

#Root: {
	parameters?: common.#Env
	chains!:     state.#Chains
	targets!:    target.#Targets
}
