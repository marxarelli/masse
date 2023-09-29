package state

#Copy: {
	copy!:       string | [string, ...string]
	from:        #Chain | *null
	destination: string

	options?: [...#CopyOption]
}
