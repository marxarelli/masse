// syntax=marxarelli/masse:v1.9.0
package main

import (
	"github.com/marxarelli/masse"
)

masse.Config

// Some top level fields to keep things DRY
buildDir: "/usr/src"
goBuildCache: "/var/cache/go"
goModCache: "/go/pkg/mod/cache"

// Our build chains, distinct parts of the overall build
chains: {
	// project is the filesystem from the main context provided by the client
	project: mainContext: true

	// builder starts with an image with the Go toolchain and sets up a build
	// directory and GOCACHE location (which we'll declare a cache for later on)
	builder: [
		{ image: "golang:1.23" },
		{ with: directory: buildDir },
		{ with: env: GOCACHE: goBuildCache },
		{ with: env: CGO_ENABLED: "0" },
	]

	// modules uses the builder environment to download Go modules according to
	// the project's go.mod and go.sum files
	//
	// It uses a `cache` mount option for `/go/pkg/mod/cache` so that subsequent
	// builds don't have to download the same modules again and again. This can
	// speed things up tremendously.
	//
	// The `diff` operation creates a resulting filesystem for this chain that
	// contains only the files created by `sh: "go mod download"`. None of the
	// files from the builder chain nor the `go.mod` and `go.sum` files are
	// present in the result.
	modules: [
		{ extend: "builder" },
		{ copy: ["go.mod", "go.sum"], from: "project" },
		{ diff: {
			sh: "go mod download"
			options: cache: goModCache
		} },
	]

	// binaries uses both the builder and modules chains (by merging their
	// filesystems) and runs `go build` to build our binary
	//
	// It uses a `cache` mount option for the `GOCACHE` directory to speed up
	// subsequent builds.
	//
	// It uses a `mount` option to mount the filesystem of the project chain at
	// the build directory. This is slightly more efficient than merging project
	// in with builder and modules, and has the caveat that files created under
	// that directory are discarded, which is why we do `go build -o ` to a
	// location outside the build directory.
	binaries: [
		{ merge: ["builder", "modules"] },
		{ sh: "go build -o /cowsayd", options: [
			{ cache: goBuildCache },
			{ mount: buildDir, from: "project" },
		] },
	]

	// cowsayd is a minimal image that contains only the `/cowsayd`
	// binary
	cowsayd: [
		{ scratch: true },
		{ copy: "/cowsayd", from: "binaries" },
	]
}

// Our build targets
targets: {
	// service is the name of our target
	service: {
		// build specifies the chain to build
		//
		// The result of this build chain is used as the image filesystem.
		build: "cowsayd"

		// runtime provides default configuration to be used when the image is run
		//
		// In constrast to Dockerfile, Massé keeps this separate from build
		// operations to clarify its function as default runtime configuration
		// unrelated to the build process.
		runtime: {
			entrypoint: ["/cowsayd"]
			user: "0"
		}
	}
}
