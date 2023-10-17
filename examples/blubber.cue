import (
	"wikimedia.org/dduvall/phyton/schema/apt"
)

parameters: {
	REPO_REMOTE: "https://gitlab.wikimedia.org/repos/releng/blubber.git"
	REPO_REF: "refs/heads/main"
}

chains: {
	repo: [
		{ git: parameters.REPO_REMOTE
			ref: parameters.REPO_REF },
	]

	go: [
		{ image: "docker-registry.wikimedia.org/golang1.19:1.19-1-20230730" },
		{ with: directory: "/src" },
	]

	tools: [
		{ extend: "go" },
		{ diff: [ apt.install & { #packages: [ "gcc", "git", "make" ] } ] },
	]

	modules: [
		{ extend: "go" },
		{ link: [ "go.mod", "go.sum" ], from: "repo" },
		{ diff: [ { run: "go mod download" } ] },
	]

	binaries: [
		{ extend: "go" },
		{ merge: ["tools", "modules"] },
		{	link: ".", from: "repo" },
		{ run: "make clean blubber-buildkit"
			options: [ { cache: "/root/.cache/go-build", access: "locked" } ] },
	]

	frontend: [
		{ copy: "blubber-buildkit", from: "binaries" }
	]
}

layouts: {
	frontend: {
		authors: [
			{ name: "Dan Duvall"
				email: "dduvall@wikimedia.org"
				keys: ["ssh-ed25519 ..."] }
		]
		comprises: ["frontend"]
		platforms: ["linux/amd64"]
	}
}
