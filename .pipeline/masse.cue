// syntax=marxarelli/masse:v1.2.0
package main

import (
	"github.com/marxarelli/masse/config"
	"github.com/marxarelli/masse-go/go"
)

config.Root

parameters: {
	goImage: string | *"docker-registry.wikimedia.org/golang1.22:1.22-20250316"
	version: string | *"v0.0.0"
}

chains: {
	projectFiles: [
		{ local: "context", options: exclude: [".git"] },
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
		go.build & { #packages: "./cmd/massed" },
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
