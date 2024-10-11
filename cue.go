package main

import (
	"fmt"
	"log"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

func main() {
	const file = `
package main

import (
	"list"
)

#State: {
	#Source | #Op
}

#Run: {
	run!: string
	arguments?: [string, ...string]
	options?: #RunOption | [#RunOption, ...#RunOption]
	if options != _|_ {
		optionsValue: list.FlattenN([options], 1)
	}
}

#RunOption: {
	#CacheMount | #SourceMount |
	#TmpFSMount | #ReadOnly
}

#CacheMount: {
	cache!: string
	access: *"shared" | "private" | "locked"
}

#SourceMount: {
	mount!: string
	from:   null | #Chain
	source: string | *"/"
}

#TmpFSMount: {
	tmpfs!: string
	size:   uint64
}

#ReadOnly: {
	readOnly!: *true | false
}

#Source: {
	#Scratch | #Image
}

#Scratch: {
	scratch: true
}

#Image: {
	image!: string
	inherit: bool | *true
}

#LayerLimit: {
	layerLimit!: uint32
}

#Op: {
	#Run
}

#Chain: [#Source, ...#Op]

#ChainRefPattern: "^[a-zA-Z_-][a-zA-Z0-9_-]*$"

#ChainRef: =~#ChainRefPattern

#Chains: {
	[=~#ChainRefPattern]: #Chain
}

chains: #Chains

chains: {
	build: [
		{ image: "debian:stable" },
		{ run: "echo hi >> /greeting" },
	],
	greet: [
		{ scratch: true },
		{ run: "cat /build/greeting", options: [ { mount: "/build", from: chains.build } ] },
	]
}
`

	ctx := cuecontext.New()

	root := ctx.CompileString(file)
	if err := root.Err(); err != nil {
		log.Fatal(err)
	}

	from := root.LookupPath(cue.ParsePath("chains.greet[1].options[0].from"))
	if from.Err() != nil {
		log.Fatal(from.Err())
	}

	fmt.Printf("path: %s\n", from.Path())

	chain := cue.Dereference(from.Eval())
	fmt.Printf("chain path: %s\n", chain.Path())
}
