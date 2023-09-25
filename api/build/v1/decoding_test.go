package v1

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	protoyaml "gitlab.wikimedia.org/dduvall/protoyaml/protoyaml"
)

func TestUnmarshalYAML(t *testing.T) {
	t.Run("Image", func(t *testing.T) {
		t.Run("fields", func(t *testing.T) {
			req := require.New(t)

			op := requireOp(t,
				"image: registry.example/foo/bar",
			)

			image := op.GetImage()
			req.NotNil(image)

			req.Equal("registry.example/foo/bar", image.Ref)
		})

		t.Run("options", func(t *testing.T) {
			t.Run("platform", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"image: registry.example/foo/bar",
					"options:",
					"- platform: linux/arm64/v7",
				)

				image := op.GetImage()
				req.NotNil(image)

				options := image.GetOptions()
				req.Len(options, 1)

				platform := options[0].GetPlatform()
				req.NotNil(platform)

				req.Equal("linux", platform.Os)
				req.Equal("arm64", platform.Architecture)
				req.Equal("v7", platform.Variant)
			})
		})
	})

	t.Run("Git", func(t *testing.T) {
		req := require.New(t)

		op := requireOp(t,
			"git: https://some.test/git/repo",
			"ref: refs/change/123",
		)

		git := op.GetGit()
		req.NotNil(git)

		req.Equal("https://some.test/git/repo", git.Remote)
		req.Equal("refs/change/123", git.Ref)
	})

	t.Run("Run", func(t *testing.T) {
		t.Run("fields", func(t *testing.T) {
			req := require.New(t)

			op := requireOp(t,
				"run: make",
				"arguments: [clean, foo]",
			)

			run := op.GetRun()
			req.NotNil(run)
			req.Equal("make", run.Command)
			req.Equal([]string{"clean", "foo"}, run.Arguments)
		})

		t.Run("options", func(t *testing.T) {
			t.Run("host", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"run: make",
					"options:",
					"- host: foo.test",
					"  ip: 1.2.3.4",
				)

				run := op.GetRun()
				req.NotNil(run)

				options := run.GetOptions()

				req.Len(options, 1)

				host := options[0].GetHost()
				req.NotNil(host)
				req.Equal("foo.test", host.Name)
				req.Equal("1.2.3.4", host.Ip)
			})

			t.Run("mount", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"run: make",
					"options:",
					"- mount: /mnt/point",
					"  from: foo-ref",
					"  source: /foo/fs/path",
				)

				run := op.GetRun()
				req.NotNil(run)

				options := run.GetOptions()

				req.Len(options, 1)

				mount := options[0].GetMount()
				req.NotNil(mount)
				req.Equal("/mnt/point", mount.Target)
				req.Equal("foo-ref", mount.From.Ref)
				req.Equal("/foo/fs/path", mount.Source)
			})

			t.Run("env", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"run: make",
					"options:",
					"- env:",
					"    FOO: bar",
				)

				run := op.GetRun()
				req.NotNil(run)

				options := run.GetOptions()

				req.Len(options, 1)

				env := options[0].GetEnv()
				req.NotNil(env)
				req.Len(env.Variables, 1)
				req.Equal("FOO", env.Variables[0].Name)
				req.Equal("bar", env.Variables[0].Value)
			})

			t.Run("cache", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"run: apt-get install cowsay",
					"options:",
					"- cache: /var/cache/apt",
					"  access: locked",
				)

				run := op.GetRun()
				req.NotNil(run)

				options := run.GetOptions()

				req.Len(options, 1)

				cache := options[0].GetCache()
				req.NotNil(cache)
				req.Equal("/var/cache/apt", cache.Target)
				req.Equal(CacheAccess_CACHE_ACCESS_LOCKED, cache.Access)
			})

			t.Run("tmpfs", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"run: make",
					"options:",
					"- tmpfs: /tmp/foo",
					"  size: 100Mb",
				)

				run := op.GetRun()
				req.NotNil(run)

				options := run.GetOptions()

				req.Len(options, 1)

				tmpfs := options[0].GetTmpfs()
				req.NotNil(tmpfs)
				req.Equal("/tmp/foo", tmpfs.Target)
				req.Equal("100Mb", tmpfs.Size)
			})
		})
	})

	t.Run("Copy", func(t *testing.T) {
		t.Run("fields", func(t *testing.T) {
			req := require.New(t)

			op := requireOp(t,
				`copy: "foo/bar"`,
				"from: foo-ref",
				"destination: ./bar/destination",
			)

			copy := op.GetCopy()
			req.NotNil(copy)

			req.Equal("foo/bar", copy.Source)
			req.Equal("foo-ref", copy.From.Ref)
			req.Equal("./bar/destination", copy.Destination)
		})

		t.Run("options", func(t *testing.T) {
			t.Run("ctime", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"copy: foo",
					"options:",
					`- ctime: "2020-01-20T01:02:03Z"`,
				)

				copy := op.GetCopy()
				req.NotNil(copy)

				options := copy.GetOptions()
				req.Len(options, 1)

				ctime := options[0].GetCtime()
				req.NotNil(ctime)
				req.Equal(int64(1579482123), ctime.Seconds)
			})

			t.Run("user", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"copy: foo",
					"options:",
					"- user: bar",
				)

				copy := op.GetCopy()
				req.NotNil(copy)

				options := copy.GetOptions()
				req.Len(options, 1)

				user := options[0].GetUser()
				req.Equal("bar", user)
			})

			t.Run("uid", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"copy: foo",
					"options:",
					"- uid: 123",
				)

				copy := op.GetCopy()
				req.NotNil(copy)

				options := copy.GetOptions()
				req.Len(options, 1)

				uid := options[0].GetUid()
				req.Equal(int32(123), uid)
			})

			t.Run("group", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"copy: foo",
					"options:",
					"- group: bar",
				)

				copy := op.GetCopy()
				req.NotNil(copy)

				options := copy.GetOptions()
				req.Len(options, 1)

				group := options[0].GetGroup()
				req.Equal("bar", group)
			})

			t.Run("gid", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"copy: foo",
					"options:",
					"- gid: 123",
				)

				copy := op.GetCopy()
				req.NotNil(copy)

				options := copy.GetOptions()
				req.Len(options, 1)

				gid := options[0].GetGid()
				req.Equal(int32(123), gid)
			})

			t.Run("mode", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"copy: foo",
					"options:",
					"- mode: 0o644",
				)

				copy := op.GetCopy()
				req.NotNil(copy)

				options := copy.GetOptions()
				req.Len(options, 1)

				mode := options[0].GetMode()
				req.Equal(uint32(0o644), mode)
			})

			t.Run("follow-symlinks", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"copy: foo",
					"options:",
					"- follow-symlinks: true",
				)

				copy := op.GetCopy()
				req.NotNil(copy)

				options := copy.GetOptions()
				req.Len(options, 1)

				follow := options[0].GetFollowSymlinks()
				req.Equal(true, follow)
			})

			t.Run("copy-dir-content", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"copy: foo",
					"options:",
					"- copy-dir-content: true",
				)

				copy := op.GetCopy()
				req.NotNil(copy)

				options := copy.GetOptions()
				req.Len(options, 1)

				dirContent := options[0].GetCopyDirContent()
				req.Equal(true, dirContent)
			})

			t.Run("include", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"copy: src/",
					"options:",
					"- include: '*.c'",
				)

				copy := op.GetCopy()
				req.NotNil(copy)

				options := copy.GetOptions()
				req.Len(options, 1)

				include := options[0].GetInclude()
				req.Equal("*.c", include)
			})

			t.Run("exclude", func(t *testing.T) {
				req := require.New(t)

				op := requireOp(t,
					"copy: src/",
					"options:",
					"- exclude: '*.js'",
				)

				copy := op.GetCopy()
				req.NotNil(copy)

				options := copy.GetOptions()
				req.Len(options, 1)

				exclude := options[0].GetExclude()
				req.Equal("*.js", exclude)
			})
		})
	})
}

func requireOp(t *testing.T, lines ...string) *Op {
	t.Helper()
	op := &Op{}
	err := protoyaml.Unmarshal([]byte(strings.Join(lines, "\n")+"\n"), op)
	require.NoError(t, err)
	return op
}
