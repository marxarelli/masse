# Massé

![masse logo](./assets/masse-256.png)

Massé is [BuildKit][buildkit] frontend that allows users to express complex
container image build graphs in [CUE][cue]. It aims to:

 1. Give users compact yet powerful declarative constructs to express how
    their container filesystems should be created, composed, and packaged.
 2. Provide composable build primitives based on (nearly) all operations
    and options in the BuildKit LLB API.
 3. Provide an API for users to define new build constructs (macros, e.g.) on
    top of build primitives.
 4. Provide a policy API for operators to assert arbitrary conditions on the
    [Low Level Build][llb] instructions prior to actual building (e.g. base
    images must come from certain registries).
 5. Formally separate container filesystem creation from image configuration.
 6. Give users a simple API for composing manifests from built filesystems and
    configuration.
 7. Support creation of manifests that contain only arbitrary meta data (e.g.
    Software Bill of Materials (SBoMs) or [attestations][in-toto-spec]).
 8. Support dynamic generation of build instructions, configuration, and
    manifest definitions via a pattern of intermediate solving. (TODO: wth
    does this mean)

## Based on CUE

The schema and user defined configuration will be written in [CUE][cue], an
"open-source data validation language and inference engine with its roots in
logic programming."

As a JSON superset, CUE has nearly the same compactness of YAML but is a very
powerful language for schema definition, validation/constraints, and user
facing configuration. With CUE we can support composable user defined types
and module imports. It's constructs are rich and coherent.

See [schema/apt/install.cue](./schema/apt/install.cue) for an example of what
an `apt install` definition looks like.

## Example config

To skip straight to what a build configuration looks like in Massé, see
Massé's own [.pipeline/masse.cue](./.pipeline/masse.cue) file which can be
used to build the BuildKit frontend image.

```console
$ docker buildx build -f .pipeline/masse.cue --target gateway .
[+] Building 27.2s (9/9) FINISHED                              remote:buildkitd
 => [internal] load build definition from masse.cue                        0.1s
 => => transferring dockerfile: 1.44kB                                     0.0s
 => resolve image config for docker-image://docker.io/marxarelli/masse:v0  1.2s
 => [auth] marxarelli/masse:pull token for registry-1.docker.io            0.0s
 => docker-image://docker.io/marxarelli/masse:v0.0.1@sha256:44866a078b236  0.8s
 => => resolve docker.io/marxarelli/masse:v0.0.1@sha256:44866a078b236ed71  0.0s
 => => sha256:bb42f0b58b5d50bdd8a20d6ace717fbf9480c3b 122.68kB / 122.68kB  0.1s
 => => sha256:a70aa5a05312867fdd6fa9686c9b4a67769d20410 16.30MB / 16.30MB  0.5s
 => => extracting sha256:a70aa5a05312867fdd6fa9686c9b4a67769d20410e514a3b  0.2s
 => => extracting sha256:bb42f0b58b5d50bdd8a20d6ace717fbf9480c3bbf0d51205  0.0s
 => local://context                                                        0.5s
 => => transferring context: 74.90MB                                       0.4s
 => docker-image://docker-registry.wikimedia.org/golang1.21:1.21-1-202311  9.7s
 => => resolve docker-registry.wikimedia.org/golang1.21:1.21-1-20231126    0.3s
 => => sha256:e11c570947367c7c5e5adb625fb56cfd3e77c35 174.42MB / 174.42MB  6.1ss
 => => sha256:c492791ecd0ad500c25716f68f37fd5c0e995f0ef 40.66MB / 40.66MB  2.9s
 => => extracting sha256:c492791ecd0ad500c25716f68f37fd5c0e995f0efbdaafad  0.7s
 => => extracting sha256:e11c570947367c7c5e5adb625fb56cfd3e77c3553ddc26b5  3.2ss
 => 📋 masse source                                                         0.4s
 => 🏗️ build `./cmd/massed`                                               13.9s 
 => 📦 package masse gateway w/ CA certificates                             0.1s
```

## Build chains

Build processes are defined as independent chains of the overall build graph.
This is meant to strike a balance between flexibility in graph node definition
while improving readability/reasoning of the overall build graph.

```cue
chains: {
  repo: [
    { git: "https://my.example/repo.git" },
  ]

  source: [
    { scratch: true },
    { copy: ".", from: "repo" },
  ]

  toolchain: [
    { image: "docker-registry.wikimedia.org/golang1.19:1.19-1-20230730" },
    { apt.install & { #packages: [ "gcc", "git", "make" ] } },
  ]

  binaries: [
    { merge: ["source", "toolchain"] },
    { with: directory: "/src" },
    { copy: ".", from: "source" },
    { run: "make" },
  ]

  application: [
    { scratch: true },
    { copy: "my-compiled-program", from: "binaries" },
  ]
}
```

As you can see with `{ merge ["source", "toolchain"] }` and `{ copy: ".",
from: "source"}` above, dependency chains are referenced by name. Chain
references are resolved during compilation.

## Macros

We can combine CUE's [definition][cue-defs] and [embedding][cue-embeds]
constructs to support a standard library and user-defined macros.

For example, the typical `apt install` pattern that's repeated in so many
Dockerfiles across the internet can be achieved with the following definition.

```cue
package apt

#PackageName: "[a-z0-9][a-z0-9+\\.\\-]+"
#VersionSpec: "(?:[0-9]+:)?[0-9]+[a-zA-Z0-9\\.\\+\\-~]*"
#ReleaseName: "[a-zA-Z](?:[a-zA-Z0-9\\-]*[a-zA-Z0-9]+)?"

#Package: =~ "^\(#PackageName)(?:=\(#VersionSpec)|/\(#ReleaseName))?$"

install: {
  #packages: [#Package, ...#Package]
  {
    sh: "apt-get install -y"
    arguments: #packages
    options: [
      { env: { "DEBIAN_FRONTEND": "noninteractive" } },
      { cache: "/var/lib/apt", access: "locked" },
      { cache: "/var/cache/apt", access: "locked" },
    ]
  }
}
```

The macro can define its parameter as a CUE definition, provide validation
constraints. In CUE terminology each package name must "unify" with the regex
constrained string.

```cue
#Package: =~ "^\(#PackageName)(?:=\(#VersionSpec)|/\(#ReleaseName))?$"
```

The macro can "return" its resulting build operations (in this case a single
`{ sh: ... }`) using an embed. The use of this shared macro is simply a CUE
unification with the definition.

```cue
import (
  "wikimedia.org/dduvall/masse/schema/apt"
)

chains:
  go: [
    { image: "docker-registry.wikimedia.org/golang1.19:1.19-1-20230730" },
    apt.install & { #packages: [ "gcc", "git", "make" ] },
  ]
```

## Modules

Massé supports CUE module dependencies and will automatically load modules
from remote OCI registries when reading in your config.

Dependencies are declared via a standard [CUE module file][cue-mod-file]
(`cue.mod/module.cue`) relative to your config file.

```
module: "project.example"
language: {
	version: "v0.13.0"
}
deps: "github.com/marxarelli/masse-go@v1": v: "v1.0.3"
```

You can import packages from these dependencies in your config, and Massé will
download the module for you at build time.

```
package main

import (
  "github.com/marxarelli/masse-go/go"
)

chains: {
  project: [
    { local: "context", options: exclude: [".git"] },
  ]

  build [
    go.image & { #version: "1.24" },
    go.mod.download & { #from: "project" },
    { copy: ".", from: "project" },
    go.build & { #packages: "./cmd/app" },
  ]
}
```

By default, Massé looks for modules in the `registry.cue.works` OCI registry.
To change this behavior, you can pass a `CUE_REGISTRY` build argument to
`docker buildx` to instruct the CUE registry client to look in a different
registry for your modules.

```
$ docker buildx build --build-arg CUE_REGISTRY=registry.example/cuemodules ...
```

See the upstream [documentation on CUE modules][cue-mod-registry-env] for
details on how to use this environment variable.

### Registry authentication

Some registries, including the default `registry.cue.works`, require
authentication even for fetching.

To use such registries, you must provide credentials (typically in the form of
a bearer token) via a BuildKit secret as well as a build argument telling
Massé what the name of the secret is.

```
$ cue login
$ REGISTRY_AUTH="$(jq -r '.registries."registry.cue.works" | (.token_type + " " + .access_token)' ~/.config/cue/logins.json)" \
  docker buildx build \
    --build-arg CUE_REGISTRY_AUTH_SECRET.registry.cue.works=REGISTRY_AUTH \
    --secret id=REGISTRY_AUTH \
    -f examples/modules/masse.cue examples/modules
```

The above command is admittedly quite cumbersome. Hopefully the CUE folks will
enable anonymous read-only access to their registry at some point and this
won't be necessary in most cases.

## TODOs

 * Organize macros into a stdlib and implement macros for most of
   [Blubber][blubber]'s higher level builder directives (npm, python, php,
   etc.).
 * Chains are currently referenced by name. While it is possible to adopt
   references by actual CUE references (e.g. `{ copy: "." from: chains.foo }`)
   this pattern encounters current performance limitations in the v2 CUE
   evaluator. The CUE folks are working on a v3 evaluator that may make this
   direct use of references possible in the future.
 * At the moment, environment variables are not substituted in command
   strings. Should they be?

## License

Massé is licensed under the GNU General Public License 3.0 or later
(GPL-3.0+). See the LICENSE file for more details.

[buildkit]: https://docs.docker.com/build/buildkit/
[llb]: https://docs.docker.com/build/buildkit/#llb
[in-toto-spec]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md
[oci]: https://github.com/opencontainers/image-spec
[frontend]: https://docs.docker.com/build/dockerfile/frontend/
[blubber]: https://gitlab.wikimedia.org/repos/releng/blubber
[cue]: https://cuelang.org
[cue-defs]: https://cuelang.org/docs/references/spec/#definitions-and-hidden-fields
[cue-embeds]: https://cuelang.org/docs/references/spec/#embedding
[cue-mod-file]: https://cuelang.org/docs/reference/modules/#cue-mod-file
[cue-mod-registry-env]: https://cuelang.org/docs/reference/modules/#cue-registry-env
