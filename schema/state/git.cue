package state

// Git uses the working directory of a given clone Git repo/ref as the initial
// filesystem for a build chain.
//
// Example:
//   chains: {
//     repo: [
//       {
//         git: "https://an.example/repo.git"
//         ref: "refs/tags/v1.2"
//         options: keepGitDir: true
//       }
//     ]
//   }
#Git: {

	// git is the URL of the Git repo to clone.
	git!:     string

	// ref is the Git branch/tag/ref to checkout.
	ref:      string | *"refs/heads/main"

	// options is one or more Git options.
	options?: #GitOption | [#GitOption, ...#GitOption]
}

#GitOption: {
	#KeepGitDir | #Constraint
}

#KeepGitDir: {

	// keepGitDir is whether the .git directory should be retained after the
	// working directory for the ref is created.
	keepGitDir!: *true | false
}
