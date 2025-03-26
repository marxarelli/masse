package target

import (
	"github.com/marxarelli/masse/common"
)

#SBOM: {
	generator: string | *""
	parameters: common.#Env
	scan: [...string] | "all" | *null
}
