package state

#Run: {
	run!:       string
	arguments?: string | [string, ...string]
	options?:   #RunOption | [#RunOption, ...#RunOption]
}
