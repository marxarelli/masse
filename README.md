# Massé

![masse logo](./assets/masse-256.png)

Masse is an extensible new [BuildKit][buildkit] frontend that allows users to
express complex container image build graphs in [CUE][cue]. It aims to:

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

See [schema/apt/macros.cue](./schema/apt/macros.cue) for an example of what an
`apt install` macro looks like.

## Example config

To skip straight to what a build configuration looks like in Masse, see
[examples/blubber.cue](./examples/blubber.cue). This is a port of Blubber's
own `blubber.yaml` to Masse.

## Build chains

Build processes are defined as independent chains of the overall build graph.
This is meant to strike a balance between flexibility in graph node definition
while improving readability/reasoning of the overall build graph.

```cue
chains: {
  repo: [
    { git: "https://my.example/repo.git" },
  ]

  toolchain: [
    { image: "docker-registry.wikimedia.org/golang1.19:1.19-1-20230730" },
    { apt.install & { #packages: [ "gcc", "git", "make" ] } },
  ]

  binaries: [
    { extend: "toolchain" },
    { with: directory: "/src" },
    { link: ".", from: "repo" },
    { run: "make" },
  ]

  application: [
    { copy: "my-compiled-program", from: "binaries" },
  ]
}
```

As you can see with `{ extend: "toolchain" }` and `{ link: ".", from: "repo"}`
above, dependency chains are referenced by name. Chain references are resolved
when the internal DAG is constructed. Cycles are also detected/prevented
during internal DAG construction.

## Macros

We can combine CUE's [definition][cuedefs] and [embedding][cueembeds]
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
    run: "apt-get install -y"
    arguments: #packages
    options: [
      { env: { "DEBIAN_FRONTEND": "noninteractive" } },
      { cache: "/var/lib/apt", access: "locked" },
      { cache: "/var/cache/apt", access: "locked" },
    ]
  }
}
```

As you can see, the macro can define its parameter as a CUE definition,
provide validation constraints. In CUE terminology each package name must
"unify" with the regex constrained string.

```cue
#Package: =~ "^\(#PackageName)(?:=\(#VersionSpec)|/\(#ReleaseName))?$"
```

The macro can "return" its resulting build operations (in this case a single
`{ run: ... }`) using an embed. The use of this shared macro is simply a CUE
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

Note that the CUE folks are working on adding syntactic sugar for this
"function" pattern. The expression could be as simple as the following in the
future.

```
apt.install(["gcc","git", "make"])
```

## TODOs

Many many things, including:

 * The `layout` specification needs attention. It should likely take the same
   approach as the `chains` specification and provide primitives that map onto
   [OCI][oci] specifications to allow for maximum flexibility in how resulting
   manifests are constructed. Perhaps the section should even be renamed
   `manifests`?
 * A BuildKit [frontend][frontend] (Dockerfile syntax) should be implemented
   soon to allow people to test this out with standard Docker tooling.
 * The `buildctl` output is wonky with emojis. It seems like width is not
   being computed correctly. This is likely an upstream bug.
 * Organize macros into a stdlib and implement macros for most of
   [Blubber][blubber]'s higher level builder directives (npm, python, php,
   etc.).
 * At the moment, environment variables are not substituted in command
   strings. Should they be?

## License

Masse is licensed under the GNU General Public License 3.0 or later
(GPL-3.0+). See the LICENSE file for more details.

[buildkit]: https://docs.docker.com/build/buildkit/
[llb]: https://docs.docker.com/build/buildkit/#llb
[in-toto-spec]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md
[cue]: https://cuelang.org
[cuedefs]: https://cuelang.org/docs/references/spec/#definitions-and-hidden-fields
[cueembeds]: https://cuelang.org/docs/references/spec/#embedding
[oci]: https://github.com/opencontainers/image-spec
[frontend]: https://docs.docker.com/build/dockerfile/frontend/
[blubber]: https://gitlab.wikimedia.org/repos/releng/blubber
