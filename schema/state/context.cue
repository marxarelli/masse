package state

#Context: {
	context!: string
	options?: #LocalOption | [#LocalOption, ...#LocalOption]
}

#MainContext: {
	mainContext!: true
	options?: #LocalOption | [#LocalOption, ...#LocalOption]
}
