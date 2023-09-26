// Core definitions
//
#Git: {
  git!: string
  ref: string | *"refs/heads/main"
}

#ImageSource: {
  image!: string
}

#Run: {
  run!: string
  arguments: [...string] | *[]
  options: [...#RunOption] | *[]
}

#RunOption: {
  #CacheMountOption | #MountOption | #EnvOption
}

#CacheMountOption: {
  cache!: string
  access: *"shared" | "private" | "locked"
}

#MountOption: {
  mount!: string
  from: #Chain
  source: string | *"/"
}

#EnvOption: {
  env!: #Env
}

#Link: {
  link!: [...string]
  from?: #Chain
  destination?: string
}

#Diff: {
  diff: #Chain
}

#Merge: {
  merge: [...#Chain]
}

#State: {
  #Git | #ImageSource | #Run | #Link | #Diff | #Merge | #StateOptions
}

#StateOptions: {
  with: [...#StateOption]
}

#StateOption: {
  #WorkingDirectory | #EnvOption
}

#WorkingDirectory: {
  working_directory: string
}

#Chain: [...#State]

#Image: {
  comprises: [...#Chain]
  authors: [...#Author]
  os: string | *"linux"
  architecture: string | *"amd64"
  configuration: #ImageConfig
}

#Author: {
  name: string
  email: string
  keys: [...#PubKey]
}

#PubKey: string

#ImageConfig: {
  user: string
  exposed_ports: {}
  environment: #Env
  entrypoint: [...string]
  default_arguments: [...string]
  working_directory: string
  labels: [...#NameValue]
  stop_signal: string
}

#Env: {
  [=~"^[a-zA-Z_][a-zA-Z0-9_]*$"]: string
}

#NameValue: {
  name: string
  value: string
}

#TargetMap: {
  [=~"."]: #Chain
}

#ImageMap: {
  [=~"."]: #Image
}

// Constraints of the root configuration fields
parameters: #Env
targets: #TargetMap
images: #ImageMap

// Macros
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

// Config start
parameters: {
  REPO_REMOTE: string | *"https://gitlab.wikimedia.org/repos/releng/blubber"
  REPO_REF: string | *"refs/heads/main"
}

targets: {
  repo: [
    {
      git: parameters.REPO_REMOTE
      ref: parameters.REPO_REF
    },
  ]

  go: [
    {
      image: "docker-registry.wikimedia.org/golang1.19:1.19-1-20230730"
    },
  ]

  build_tools: [
    {
      merge: [targets.go]
    },
    #AptInstall & {
      #packages: ["gcc", "git", "make"]
    },
  ]

  module: [
    {
      merge: [targets.go]
    },
    {
      link: ["go.mod", "go.sum"]
      from: targets.repo
    },
    {
      diff: [
        { run: "go mod download" },
      ]
    },
  ]

  binaries: [
    {
      merge: [targets.go, targets.build_tools, targets.modules]
    },
    {
      with: [
        {
          working_directory: "/src"
        }
      ]
    },
    {
      link: ["."]
      from: targets.repo
    },
    {
      run: "make clean blubber-buildkit",
      options: [
        {
          cache: "/root/.cache/go-build"
          access: "locked"
        }
      ]
    }
  ]

  frontend: [
    {
      link: ["/src/blubber-buildkit"]
      from: targets.binaries
      destination: "/blubber-buildkit"
    }
  ]
}

layouts: {
  frontend: {
    authors: [
      {
        name: "Dan Duvall"
        email: "dduvall@wikimedia.org"
        keys: ["ssh-ed25519 ..."]
      }
    ]
    comprises: [
      targets.frontend
    ]
  }
}
