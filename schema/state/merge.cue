package state

// Merge creates an aggregate filesystem from a number of build chains.
#Merge: {

	// merge is one or more build chain names.
	merge!: #ChainRef | [#ChainRef, ...#ChainRef]
}
