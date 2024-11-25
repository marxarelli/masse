package state

#Image: {
	image!:   string
	inherit:  bool | *true
	options?: #ImageOption | [#ImageOption, ...#ImageOption]
}

#ImageOption: {
	#LayerLimit | #Constraint
}

#LayerLimit: {
	layerLimit!: uint32
}
