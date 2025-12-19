// syntax=marxarelli/masse:v1.10.0
package main

import (
	"github.com/marxarelli/masse"
	"github.com/marxarelli/masse-go/go"
)

masse.Config

parameters: {
	goImage: string | *"docker-registry.wikimedia.org/golang1.24:1.24-20251214"
	tag: string | *"v0.0.0"
}

chains: {
	projectFiles: [
		{
			mainContext: true,
			options: exclude: [".git", "bin", "build", "massed", "masse"]
		}
	]

	goImage: [
		{ image: parameters.goImage },
	]

	modules: [
		{ extend: "goImage" },
		go.mod.download & { #from: "projectFiles" },
	]

	build: [
		{ merge: ["goImage", "modules"] },
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
		go.build.sh & { #command: "make TAG=\(parameters.tag)", #arguments: "bin/massed" },
	]

	gateway: [
		{ scratch: true },
		{
			file: [
				{
					copy: "/src/bin/massed"
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
