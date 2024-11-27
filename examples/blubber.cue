package main

import (
	"wikimedia.org/dduvall/masse/apt"
)

parameters: {
	REPO_REMOTE: string | *"https://gitlab.wikimedia.org/repos/releng/blubber.git"
	REPO_REF: string | *"refs/heads/main"
}

chains: {
	repo: [
		{ git: parameters.REPO_REMOTE
			ref: parameters.REPO_REF },
	]

	go: [
		{ image: "docker-registry.wikimedia.org/golang1.22:1.22-20241124" },
		{ with: directory: "/src" },
	]

	tools: [
		{ extend: "go" },
		{ diff: [ apt.install & { #packages: [ "gcc", "git", "make" ] } ] },
	]

	modules: [
		{ extend: "go" },
		{ file: [
			{ copy: "go.mod", from: "repo" },
			{ copy: "go.sum", from: "repo" },
		] },
		{ diff: [ { run: "go mod download" } ] },
	]

	binaries: [
		{ extend: "go" },
		{ merge: ["tools", "modules"] },
		{ file: { copy: ".", from: "repo" } },
		{ run: "make clean blubber-buildkit"
			options: [ { cache: "/root/.cache/go-build", access: "locked" } ] },
	]
}

targets: {
	frontend: {
		build: [
			{ scratch: true },
			{ file: { copy: "blubber-buildkit", from: "binaries" } },
		],
		platforms: ["linux/amd64", "linux/arm64"]
		runtime: {
			user: "nobody"
			entrypoint: ["/blubber-buildkit"]
		}
	}
}
