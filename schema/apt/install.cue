package apt

import (
	"strings"
)

#PackageName: "[a-z0-9][a-z0-9+\\.\\-]+"
#VersionSpec: "(?:[0-9]+:)?[0-9]+[a-zA-Z0-9\\.\\+\\-~]*"
#ReleaseName: "[a-zA-Z](?:[a-zA-Z0-9\\-]*[a-zA-Z0-9]+)?"

#Package: =~ "^\(#PackageName)(?:=\(#VersionSpec)|/\(#ReleaseName))?$"

install: {
	#packages: [#Package, ...#Package]
	{
		sh: "apt-get update && apt-get install -y"
		arguments: #packages
		options: [
			{ env: { "DEBIAN_FRONTEND": "noninteractive" } },
			{ cache: "/var/lib/apt", access: "locked" },
			{ cache: "/var/cache/apt", access: "locked" },
			{ customName: "📦 installing APT packages (" + strings.Join(#packages, " ") + ")" },
		]
	}
}
