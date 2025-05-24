package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/masse/common"
	"gitlab.wikimedia.org/dduvall/masse/load"
)

func TestLoad(t *testing.T) {
	req := require.New(t)
	dir := t.TempDir()

	root, err := Load(
		filepath.Join(dir, "masse.cue"),
		[]byte(`
package main

import (
	"github.com/marxarelli/masse"
	"github.com/marxarelli/masse/apt"
)

masse.Config

parameters: {
	repo: string | *"https://gitlab.wikimedia.org/repos/releng/blubber.git"
	ref: string | *"refs/heads/main"
}

chains: {
	repo: [
		{ git: parameters.repo, ref: parameters.ref },
	]

	go: [
		{ image: "docker-registry.wikimedia.org/golang1.19:1.19-1-20230730" },
		{ with: directory: "/src" },
	]

	tools: [
		{ extend: go },
		{ diff: apt.install & { #packages: [ "gcc", "git", "make" ] } },
	]

	modules: [
		{ extend: go },
		{ file: [
				{ copy: "go.mod", from: repo },
				{ copy: "go.sum", from: repo },
		] },
		{ diff: [ { sh: "go mod download" } ] },
	]

	binaries: [
		{ extend: go },
		{ merge: [tools, modules] },
		{	file: { copy: ".", from: repo } },
		{ sh: "make clean blubber-buildkit"
			options: [ { cache: "/root/.cache/go-build", access: "locked" } ] },
	]

	frontend: [
		{ scratch: true },
		{ file: { copy: "blubber-buildkit", from: binaries } },
	]
}

targets: {
	#default: {
		platforms: ["linux/amd64", "linux/arm64"]
		runtime: {
			user: "nobody"
			entrypoint: ["/default"]
		}
		attestations: sbom: generator: "an.example/sbom/generator:v1.0"
	}

	frontend: {
		build: chains.frontend
		runtime: {
			entrypoint: ["/blubber-buildkit"]
		}
		attestations: sbom: generator: "an.example/sbom/generator:v2.0"
	}
}
		`),
		map[string]string{
			"ref": `"refs/tags/v1.0"`,
		},
		load.WithDefaultEmbeddedModFile(),
	)

	req.NoError(err)
	req.Contains(root.Targets, "frontend")

	frontend := root.Targets["frontend"]
	req.Equal(
		[]common.Platform{
			{OS: "linux", Architecture: "amd64"},
			{OS: "linux", Architecture: "arm64"},
		},
		frontend.Platforms,
	)
	req.Equal(
		"nobody",
		frontend.Runtime.User,
	)

	sbom := frontend.Attestations.SBOM
	req.Equal(
		"an.example/sbom/generator:v2.0",
		sbom.Generator,
	)
}

func BenchmarkLoadWithStringReferences(b *testing.B) {
	dir := b.TempDir()
	for i := 0; i < b.N; i++ {

		Load(
			filepath.Join(dir, "masse.cue"),
			[]byte(`
package main

import (
	"github.com/marxarelli/masse"
)

masse.Config

buildDir: "/usr/src"

chains: {
	project: mainContext: true

	builder: [
		{ image: "golang:1.23" },
		{ with: directory: buildDir },
	]

	modules: [
		{ extend: "builder" },
		{ copy: ["go.mod", "go.sum"], from: "project" },
		{ diff: sh: "go mod download" },
	]

	binaries: [
		{ merge: ["builder", "modules"] },
		{ sh: "go build -o /cowsayd", options: [
			{ mount: buildDir, from: "project" },
		] },
	]

	cowsayd: [
		{ scratch: true },
		{ copy: "/cowsayd", from: "binaries" },
	]
}

targets: {
	service: {
		build: "cowsayd"
	}
}
			`),
			map[string]string{
				"ref": `"refs/tags/v1.0"`,
			},
			load.WithDefaultEmbeddedModFile(),
		)
	}
}

func BenchmarkLoadWithCUEReferences(b *testing.B) {
	dir := b.TempDir()
	for i := 0; i < b.N; i++ {

		Load(
			filepath.Join(dir, "masse.cue"),
			[]byte(`
package main

import (
	"github.com/marxarelli/masse"
)

masse.Config

buildDir: "/usr/src"

chains: {
	project: mainContext: true

	builder: [
		{ image: "golang:1.23" },
		{ with: directory: buildDir },
	]

	modules: [
		{ extend: "builder" },
		{ copy: ["go.mod", "go.sum"], from: project },
		{ diff: sh: "go mod download" },
	]

	binaries: [
		{ merge: ["builder", "modules"] },
		{ sh: "go build -o /cowsayd", options: [
			{ mount: buildDir, from: project },
		] },
	]

	cowsayd: [
		{ scratch: true },
		{ copy: "/cowsayd", from: binaries },
	]
}

targets: {
	service: {
		build: chains.cowsayd
	}
}
			`),
			map[string]string{
				"ref": `"refs/tags/v1.0"`,
			},
			load.WithDefaultEmbeddedModFile(),
		)
	}
}
