package layout

import (
	"wikimedia.org/dduvall/phyton/common"
)

#Image: {
	comprises: [...#Chain]
	authors: [...#Author]
	platform:      common.#Platform
	configuration: #ImageConfig
}
