package state

// Extend uses the filesystem result of another build chain as a starting point for a
// new chain.
//
// Example:
//   chains: {
//     repo: [
//       { git: "https://an.example/repo.git" },
//     ]
//     build: [
//       { extends: "repo" },
//       { run: "make" }
//     ]
//   }
#Extend: {

	// extend references another build chain by name.
	extend!: #ChainRef
}
