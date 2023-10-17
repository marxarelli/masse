package state

#Image: {
	image!: string

	options?: [...#ImageOption]
}

#ImageOption: {
	#LayerLimit | #Constraint
}

#LayerLimit: {
	layerLimit!: uint32
}
