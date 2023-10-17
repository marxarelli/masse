package state

import (
	"list"
)

#Git: {
	git!: string
	ref:  string | *"refs/heads/main"
	options?: #GitOption | [#GitOption, ...#GitOption]
	if options != _|_ {
		optionsValue: list.FlattenN([options], 1)
	}
}

#GitOption: {
	#KeepGitDir | #Constraint
}

#KeepGitDir: {
	keepGitDir!: *true | false
}
