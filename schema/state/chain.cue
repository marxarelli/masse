package state

// Chain is a list that starts with a build source and follows with a list
// of serial operations that result in an image filesystem.
#Chain: #Source | [#Source, ...#Op]

#ChainRefPattern: "^[a-zA-Z_-][a-zA-Z0-9_-]*$"

#ChainRef: =~#ChainRefPattern
#ChainStringRef: =~#ChainRefPattern

// Chains is a number of named composable build chains that are used to
// produce an image filesystem.
#Chains: {
	[=~#ChainRefPattern]: #Chain
}
