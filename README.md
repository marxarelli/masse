# Phyton

> One of the parts which by their repetition make up a flowering plant, each
> being a single joint of a stem with its leaf or leaves; a phytomer.

– phyton. (2023) [Wiktionary](https://en.wiktionary.org/wiki/phyton).


Phyton is an extensible new [BuildKit][buildkit] frontend that allows users to
express complex container image build graphs in. It aims to:

 1. Give users compact yet powerful declarative constructs to express how
    their container filesystems should be created, composed, and packaged.
 2. Provide an API for users to define new build constructs (macros, e.g.).
 3. Provide a policy API for operators to constrain privileged operations.
 4. Maintain a lazy evaluation model by expressing all build instructions as
    [Low-Level Build (LLB)][llb].
 5. Formally separate container filesystem creation from image configuration.
 6. Give users a simple API for composing images from built filesystems and
    configuration.
 7. Bake supply chain security into the specification itself. Can the layout
    itself subsume an [in-toto assertions][in-toto-spec]?

## Known unknowns

### Specification and configuration language

Oh to choose.

#### YAML

YAML is well known. It is also notorious for its ambiguity and abuse.

By itself, it cannot satisfy the goal of provide user defined extension
(macros).  It would need something like embedded CEL and esoteric YAML
mappings for real language constructs and evaluation.

```
macros:
  apt::packages:
    - assert: # validation or guards akin to Haskell
        $packages: []
        $options: []
        where: "packages.all(p, p.matches(debianPackage))"
      invalid: "Invalid package names"
      =>: !cel |
        [
          Run {
            command: "apt-get update && apt-get install --no-recommends -y",
            args: packages,
            options: [
              env({

              }),
              Cache {
                path: "/var/cache/apt",
                access: Cache.LOCKED,
              },
            ] + runOptions(options),
          },
        ]

targets:
  go:
    - image: docker-registry.wikimedia.org/golang1.19:1.19-1-20230730

  build-tools:
    - merge: [go]
    - diff:
        - apt::packages: [gcc, git, make]

```

See [examples/blubber.yaml](./examples/blubber.yaml) for what a configuration
might look like in YAML with embedded CEL expressions for macros.

#### CEL + protobuf

CEL expressions are capable of representing a build graph on their own.
However, a large graph would contain a lot of type-declaration overhead,
hardly the compact format needed here.

CEL (or YAML + embedded CEL) would require the schema to be written in Proto
definitions to get the data structure into Go for evaluation to [LLB][llb].

The protobuf schema would need an additional validation layer as well, likely
also in CEL using something like [protoc-gen-validate][protoc-gen-validate].

#### CUE

[CUE][cue] is very interesting. It has nearly the same compactness of YAML but
is a very powerful language for schema definition, validation/constraints, and
user configuration. It's constructs are rich and coherent.

Macros, for instance, can be accomplished quite easily using embedded
definitions and hidden fields.

```cue
#Chain: [...#State]

#State: {
  #Git | #ImageSource | #Run | #Link | #Diff | #Merge | #StateOptions
}

#Diff: {
  diff: #Chain
}

// #AptInstall unifies with #Diff since its #packages field is hidden and so
// it is a valid #State
#AptInstall: #Diff & {
  #packages: [...string]
  diff: [
    {
      run: "apt-get install -y"
      arguments: #packages
      options: [
        {
          cache: "/var/lib/apt"
          access: "locked"
        },
        {
          cache: "/var/cache/apt"
          access: "locked"
        }
      ]
    },
  ]
}

targets: {
  go: [ { image: "docker-registry.wikimedia.org/golang1.19:1.19-1-20230730" } ]

  build_tools: [
    { merge: [targets.go] },
    #AptInstall & {
      #packages: ["gcc", "git", "make"]
    },
  ]
```

Note that in the above, `targets.go` is referencing directly the `go` field of
`targets`, another very powerful construct akin to YAML anchors but much
better.

CUE also supports imports, meaning users would be able to publish, share,
import macros.

## License

Phyton is licensed under the GNU General Public License 3.0 or later
(GPL-3.0+). See the LICENSE file for more details.

[buildkit]: https://docs.docker.com/build/buildkit/
[llb]: https://docs.docker.com/build/buildkit/#llb
[in-toto-spec]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md
[protoc-gen-validate]: https://github.com/bufbuild/protoc-gen-validate
[cue]: https://cuelang.org
