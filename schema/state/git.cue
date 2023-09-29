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
	keep_git_dir: *true | false
}
