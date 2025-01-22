package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

func TestLoad(t *testing.T) {
	req := require.New(t)
	root, err := Load(
		"masse.cue",
		[]byte(`
package main

import (
	"github.com/marxarelli/masse/apt"
)

parameters: {
	repo: string | *"https://gitlab.wikimedia.org/repos/releng/blubber.git"
	ref: string | *"refs/heads/main"
}

chains: {
	repo: [
		{ git: parameters.repo
			ref: parameters.ref },
	]

	go: [
		{ image: "docker-registry.wikimedia.org/golang1.19:1.19-1-20230730" },
		{ with: directory: "/src" },
	]

	tools: [
		{ extend: "go" },
		{ diff: apt.install & { #packages: [ "gcc", "git", "make" ] } },
	]

	modules: [
		{ extend: "go" },
		{ file: [
				{ copy: "go.mod", from: "repo" },
				{ copy: "go.sum", from: "repo" },
		] },
		{ diff: [ { sh: "go mod download" } ] },
	]

	binaries: [
		{ extend: "go" },
		{ merge: ["tools", "modules"] },
		{	file: { copy: ".", from: "repo" } },
		{ sh: "make clean blubber-buildkit"
			options: [ { cache: "/root/.cache/go-build", access: "locked" } ] },
	]

	frontend: [
		{ scratch: true },
		{ file: { copy: "blubber-buildkit", from: "binaries" } },
	]
}

targets: {
	#default: {
		platforms: ["linux/amd64", "linux/arm64"]
		runtime: user: "nobody"
	}

	frontend: {
		build: "frontend"
		runtime: {
			entrypoint: ["/blubber-buildkit"]
		}
	}
}
		`),
		map[string]string{
			"ref": `"refs/tags/v1.0"`,
		},
	)

	req.NoError(err)
	req.Contains(root.Targets, "frontend")

	frontend := root.Targets["frontend"]
	req.Equal(
		frontend.Platforms,
		[]common.Platform{
			{OS: "linux", Architecture: "amd64"},
			{OS: "linux", Architecture: "arm64"},
		},
	)
	req.Equal(
		frontend.Runtime.User,
		"nobody",
	)
}
