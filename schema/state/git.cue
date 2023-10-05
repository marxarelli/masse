package state

#Git: {
	git!: string
	ref:  string | *"refs/heads/main"
	options?: [...#GitOption]
}

#GitOption: {
	#KeepGitDir
}

#KeepGitDir: {
	keepGitDir!: *true | false
}
