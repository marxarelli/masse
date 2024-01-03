package config

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/masse/common"
	"gitlab.wikimedia.org/dduvall/masse/state"
	"gitlab.wikimedia.org/dduvall/masse/target"
	"gitlab.wikimedia.org/dduvall/masse/util/testdecode"
)

func TestDecodeLayout(t *testing.T) {
	tester := &testdecode.Tester{
		T: t,
		CUEImports: []string{
			"wikimedia.org/dduvall/masse/schema/apt",
			"wikimedia.org/dduvall/masse/schema/config",
		},
	}

	testdecode.Run(tester,
		"config.#Root",
		`config.#Root & {
			parameters: {
				REPO_REMOTE: "https://some.example/repo.git"
				REPO_REF: "refs/heads/main"
			}

			chains: {
				repo: [
					{ git: parameters.REPO_REMOTE
						ref: parameters.REPO_REF },
				]

				go: [
					{ image: "docker-registry.wikimedia.org/golang1.19:1.19-1-20230730" },
				]

				tools: [
					{ merge: [ "go" ] },
					apt.install & { #packages: [ "gcc", "git", "make" ] },
				]

				modules: [
					{ merge: [ "go" ] },
					{ link: [ "go.mod", "go.sum" ], from: "repo" },
					{ diff: [ { run: "go mod download" } ] },
				]

				binaries: [
					{ merge: [ "go", "tools", "modules" ] },
					{ with: [ { directory: "/src" } ] },
					{	link: ".", from: "repo" },
					{ run: "make clean blubber-buildkit"
					  options: [ { cache: "/root/.cache/go-build", access: "locked" } ] },
				]

				frontend: [
					{ copy: "/src/blubber-buildkit",
					  from: "binaries",
						destination: "/blubber-buildkit" }
				]
			}

			targets: {
				frontend: {
					build: "frontend"
					platforms: ["linux/amd64", "linux/arm64"]
					runtime: {
						user: "nobody"
						entrypoint: ["/blubber-buildkit"]
					}
				}
			}
		}`,
		Root{
			Parameters: Parameters{
				"REPO_REMOTE": "https://some.example/repo.git",
				"REPO_REF":    "refs/heads/main",
			},
			Chains: state.Chains{
				"repo": state.Chain{
					{Git: &state.Git{
						Repo: "https://some.example/repo.git",
						Ref:  "refs/heads/main",
					}},
				},
				"go": state.Chain{
					{Image: &state.Image{Ref: "docker-registry.wikimedia.org/golang1.19:1.19-1-20230730", Inherit: true}},
				},
				"tools": state.Chain{
					{Merge: &state.Merge{Merge: []state.ChainRef{"go"}}},
					{Run: &state.Run{
						Command:   "apt-get install -y",
						Arguments: []string{"gcc", "git", "make"},
						Options: []*state.RunOption{
							{Option: &state.Option{Env: &state.Env{Env: common.Env{"DEBIAN_FRONTEND": "noninteractive"}}}},
							{CacheMount: &state.CacheMount{Target: "/var/lib/apt", Access: state.CacheLocked}},
							{CacheMount: &state.CacheMount{Target: "/var/cache/apt", Access: state.CacheLocked}},
						},
					}},
				},
				"modules": state.Chain{
					{Merge: &state.Merge{Merge: []state.ChainRef{"go"}}},
					{Link: &state.Link{Source: []common.Glob{"go.mod", "go.sum"}, Destination: "./", From: "repo"}},
					{Diff: &state.Diff{Upper: state.Chain{
						{Run: &state.Run{Command: "go mod download"}},
					}}},
				},
				"binaries": state.Chain{
					{Merge: &state.Merge{Merge: []state.ChainRef{"go", "tools", "modules"}}},
					{With: &state.With{With: []*state.Option{
						{WorkingDirectory: &state.WorkingDirectory{Directory: "/src"}},
					}}},
					{Link: &state.Link{Source: []common.Glob{"."}, Destination: "./", From: "repo"}},
					{Run: &state.Run{
						Command: "make clean blubber-buildkit",
						Options: []*state.RunOption{
							{CacheMount: &state.CacheMount{Target: "/root/.cache/go-build", Access: state.CacheLocked}},
						},
					}},
				},
				"frontend": state.Chain{
					{Copy: &state.Copy{
						Source:      []common.Glob{"/src/blubber-buildkit"},
						From:        "binaries",
						Destination: "/blubber-buildkit",
					}},
				},
			},
			Targets: target.Targets{
				"frontend": &target.Target{
					Build: state.ChainRef("frontend"),
					Platforms: []common.Platform{
						{OS: "linux", Architecture: "amd64"},
						{OS: "linux", Architecture: "arm64"},
					},
					Runtime: target.Runtime{
						User:       common.User{User: "nobody"},
						Entrypoint: []string{"/blubber-buildkit"},
						Env:        common.Env{},
						Directory:  "/",
						StopSignal: "SIGTERM",
					},
					Labels: map[string]string{},
				},
			},
		},
	)
}
