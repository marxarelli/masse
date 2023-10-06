package layout

import (
	"wikimedia.org/dduvall/phyton/schema/state"
)

#Root: {
	parameters?: #Parameters
	chains!:     state.#Chains
	layouts!:    #Layouts
}
