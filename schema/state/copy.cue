package state

#Copy: {
	copy!:  string | [string, ...string]
	source: [
		if ((copy & string) != _|_) {
			[copy]
		},
		copy,
	][0]
	from:        #ChainRef
	destination: string | *"./"

	options?: [...#CopyOption]
}
