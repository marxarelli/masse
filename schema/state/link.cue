package state

#Link: {
	link!:  string | [string, ...string]
	source: [
		if ((link & string) != _|_) {
			[link]
		},
		link,
	][0]
	from:        #ChainRef
	destination: string | *"./"

	options?: [...#CopyOption]
}
