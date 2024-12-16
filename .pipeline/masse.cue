// syntax=marxarelli/masse:v0.0.1
package main

parameters: {
	goImage: string | * "docker-registry.wikimedia.org/golang1.21:1.21-1-20231126"
}

chains: {
	local: [
		{ local: "context" },
	]

	go: [
		{ image: parameters.goImage },
	]

	modules: [
		{ extend: "go" },
		{ file: [
			{ mkdir: "/src" },
			{ copy: "go.mod", destination: "/src", from: "local" },
			{ copy: "go.sum", destination: "/src", from: "local" },
		] },
		{ diff: run: "go mod download" },
	]

	build: [
		{ extend: "go" },
		{ with: directory: "/src" },
		{ with: env: CGO_ENABLED: "0" },
		{
			file: [
				{ copy: ".",  from: "local" },
			]
			options: customName: "📋 masse source"
		},
	]

	massed: [
		{ extend: "build" },
		{
			run: "go build ./cmd/massed"
			options: customName: "🏗️ build `./cmd/massed`"
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
		build: [
			{ extend: "gateway" },
		]
		runtime: {
			user: "nobody"
			entrypoint: ["/massed"]
		}
	}
}
