package state

#Git: {
	git!:     string
	ref:      string | *"refs/heads/main"
	options?: #GitOption | [#GitOption, ...#GitOption]
}

#GitOption: {
	#KeepGitDir | #Constraint
}

#KeepGitDir: {
	keepGitDir!: *true | false
}
