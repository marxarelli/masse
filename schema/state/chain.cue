package state

#Chain: [#Source, ...#Op]

#ChainRefPattern: "^[a-zA-Z_-][a-zA-Z0-9_-]*$"

#ChainRef: =~#ChainRefPattern

#Chains: {
	[=~#ChainRefPattern]: #Chain
}
