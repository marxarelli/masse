package npm

install: {
	#command: "ci" | *"install"
	#packageFile: string | *"package.json"
	#lockFile: string | *"package-lock.json"
	#environment: string | *"production"
	#cache: string | *"/var/lib/cache/npm"

	{
		ops: [
			{
				run: "npm install"
				arguments: [
					if #environment == "production" {
						"--only=production"
					}
				]
				options: [
					{ env: { NPM_CONFIG_CACHE: #cache } },
					{ cache: #cache, access: "locked" },
				]
			},
			{
				run: "npm dedupe"
				options: [
					{ env: { NPM_CONFIG_CACHE: #cache } },
					{ cache: #cache, access: "locked" },
				]
			}
		]
	}
}
