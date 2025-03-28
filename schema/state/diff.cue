package state

// Diff creates a filesystem containing only the difference between that of
// the current build chain (lower) and the same build chain with the given
// operations applied (upper).
#Diff: {

	// diff is one or more operations to apply to the build chain before
	// creating the diff filesystem.
	diff!: #Op | [#Op, ...#Op] | *null

	// from is the optional lower state against which to diff. If it is
	// omitted, the state of the current build chain is used.
	from?: string

	// options is one or more diff options.
	options?: #Constraint | [#Constraint, ...#Constraint]
}
