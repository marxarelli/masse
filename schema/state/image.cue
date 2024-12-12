package state

// Image uses the filesystem and configuration of an existing OCI image as the
// initial filesystem for a build chain.
#Image: {

	// image is a reference to a remote OCI image
	image!:   string

	// inherit is whether the image configuration should be applied to the build
	// chain.
	inherit:  bool | *true

	// options is one or more image options.
	options?: #ImageOption | [#ImageOption, ...#ImageOption]
}

#ImageOption: {
	#LayerLimit | #Constraint
}

#LayerLimit: {

	// layerLimit is the maximum number of layers allowed in a source image.
	layerLimit!: uint32
}
