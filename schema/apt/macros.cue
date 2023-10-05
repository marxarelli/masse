package apt

#Packages: [string, ...string]

#Install: {
	packages: #Packages
	out: {
		run: "apt-get install -y"
		arguments: packages
		options: [
			{ env: { "DEBIAN_FRONTEND": "noninteractive" } },
			{ cache: "/var/lib/apt", access: "locked" },
			{ cache: "/var/cache/apt", access: "locked" },
		]
	}
}
