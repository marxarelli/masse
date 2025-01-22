// syntax=marxarelli/masse:experimental
package main

parameters: {
	goImage: string | * "docker-registry.wikimedia.org/golang1.22:1.22-20250316"
}

chains: {
	projectFiles: [
		{ local: "context", options: exclude: [".git"] },
	]

	go: [
		{ image: parameters.goImage },
	]

	modules: [
		{ extend: "go" },
		{ file: [
			{ mkdir: "/src" },
			{ copy: "go.mod", destination: "/src", from: "projectFiles" },
			{ copy: "go.sum", destination: "/src", from: "projectFiles" },
		] },
		{ diff: { run: ["go", "mod", "download"], options: directory: "/src" } },
	]

	build: [
		{ merge: ["go", "modules"] },
		{ with: directory: "/src" },
		{ with: env: CGO_ENABLED: "0" },
		{
			file: [
				{ copy: ".",  from: "projectFiles" },
			]
			options: customName: "📋 masse source"
		},
	]

	massed: [
		{ extend: "build" },
		{
			run: ["go", "build", "./cmd/massed"]
			options: [
				{ customName: "🏗️ build `./cmd/massed`" },
				{ cache: "/var/cache/go", access: "locked" },
				{ env: GOCACHE: "/var/cache/go" },
			]
		},
	]

	gateway: [
		{ scratch: true },
		{
			file: [
				{
					copy: "/src/massed"
					destination: "/massed"
					from: "massed"
				},
				{
					mkdir: "/etc/ssl/certs"
					options: createParents: true
				},
				{
					copy: "/etc/ssl/certs/ca-certificates.crt"
					destination: "/etc/ssl/certs/ca-certificates.crt"
					from: "massed"
				},
			]
			options: customName: "📦 package masse gateway w/ CA certificates"
		},
	]
}

targets: {
	gateway: {
		platforms: ["linux/amd64"]
		build: "gateway"
		runtime: {
			user: "nobody"
			entrypoint: ["/massed"]
		}
	}
}
