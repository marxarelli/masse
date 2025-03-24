package config

import (
	"github.com/marxarelli/masse/common"
	"github.com/marxarelli/masse/state"
	"github.com/marxarelli/masse/target"
)

Root: {
	parameters?: common.#Env
	chains!:     state.#Chains
	targets!:    target.#Targets
}
