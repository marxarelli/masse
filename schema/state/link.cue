package state

#Link: {
	link!:       string | [string, ...string]
	from:        #Chain | *null
	destination: string

	options?: [...#CopyOption]
}
