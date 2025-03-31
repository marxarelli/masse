// syntax=marxarelli/masse:v1.7.0
package main

import (
	"github.com/marxarelli/masse"
	"github.com/marxarelli/masse-go/go"
)

masse.Config

parameters: {
	goImage: string | *"docker-registry.wikimedia.org/golang1.22:1.22-20250316"
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
		go.build & { #command: "make TAG=\(parameters.tag)", #packages: "bin/massed" },
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
