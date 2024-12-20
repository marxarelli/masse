package v1

import (
	"testing"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestFile(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/dduvall/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			c := newCompiler(nil)
			c.chainCompilers = map[string]chainCompiler{
				"local": func(_ *compiler) *chainResult {
					return &chainResult{state: llb.Local("context").Dir("/src")}
				},
			}
			return c
		}),
	)

	compile.Run("options", func(compile *testcompile.Tester) {
		compile.Test(
			"customName",
			`state.#File & { file: { copy: "./foo", from: "local" }, options: customName: "copying foo" }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				ops, _ := req.ContainsNFileOps(1)
				md := req.HasMetadata(ops[0])
				req.Contains(md.Description, "llb.customname")
				req.Equal("copying foo", md.Description["llb.customname"])
			},
		)
	})

	compile.Run("Copy", func(compile *testcompile.Tester) {
		compile.Test(
			"minimal",
			`state.#File & { file: { copy: "./foo", from: "local" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.Equal("/src/foo", copies[0].Copy.Src)
				req.Equal("/", copies[0].Copy.Dest)

				req.False(copies[0].Copy.DirCopyContents)
				req.False(copies[0].Copy.FollowSymlink)
			},
		)

		compile.Test(
			"destination",
			`state.#File & { file: { copy: "/foo/bar", from: "local", destination: "./baz" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.Equal("/foo/bar", copies[0].Copy.Src)
				req.Equal("/dest/baz", copies[0].Copy.Dest)
			},
			testcompile.WithInitialState(llb.Scratch().Dir("/dest")),
		)

		compile.Test(
			"options/ctime",
			`state.#File & { file: { copy: "./to", from: "local", destination: "./the/sound", options: ctime: "2016-04-11T17:23:07-07:00" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.Equal("/src/to", copies[0].Copy.Src)
				req.Equal("/the/sound", copies[0].Copy.Dest)
				req.Equal(int64(1460420587000000000), copies[0].Copy.Timestamp)
			},
		)

		compile.Test(
			"options/user/name",
			`state.#File & { file: { copy: "./foo", from: "local", options: user: "kim" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.Equal("kim", copies[0].Copy.Owner.User.User.(*pb.UserOpt_ByName).ByName.Name)
			},
		)

		compile.Test(
			"options/user/uid",
			`state.#File & { file: { copy: "./foo", from: "local", options: uid: 9 } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.Equal(uint32(9), copies[0].Copy.Owner.User.User.(*pb.UserOpt_ByID).ByID)
			},
		)

		compile.Test(
			"options/group/name",
			`state.#File & { file: { copy: "./foo", from: "local", options: group: "breeders" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.Equal("breeders", copies[0].Copy.Owner.Group.User.(*pb.UserOpt_ByName).ByName.Name)
			},
		)

		compile.Test(
			"options/group/gid",
			`state.#File & { file: { copy: "./foo", from: "local", options: gid: 1 } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.Equal(uint32(1), copies[0].Copy.Owner.Group.User.(*pb.UserOpt_ByID).ByID)
			},
		)

		compile.Test(
			"options/copyDirectoryContents",
			`state.#File & { file: { copy: "./foo", from: "local", options: copyDirectoryContents: true } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.True(copies[0].Copy.DirCopyContents)
			},
		)

		compile.Test(
			"options/followSymlinks",
			`state.#File & { file: { copy: "./foo", from: "local", options: followSymlinks: true } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.True(copies[0].Copy.FollowSymlink)
			},
		)

		compile.Test(
			"options/allowNotFound",
			`state.#File & { file: { copy: "./foo", from: "local", options: allowNotFound: true } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.True(copies[0].Copy.AllowEmptyWildcard)
			},
		)

		compile.Test(
			"options/wildcard",
			`state.#File & { file: { copy: "./foo/*", from: "local", options: wildcard: true } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.True(copies[0].Copy.AllowWildcard)
			},
		)

		compile.Test(
			"options/replaceExisting",
			`state.#File & { file: { copy: "./foo", destination: "./foo", from: "local", options: replaceExisting: true } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.True(copies[0].Copy.AlwaysReplaceExistingDestPaths)
			},
		)

		compile.Test(
			"options/createParents",
			`state.#File & { file: { copy: "./bar", destination: "/foo/bar", from: "local", options: createParents: true } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)
				req.True(copies[0].Copy.CreateDestPath)
			},
		)
	})

	compile.Run("Mkfile", func(compile *testcompile.Tester) {
		compile.Test(
			"minimal",
			`state.#File & { file: { mkfile: "./foo", content: 'some content' } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkfiles := req.ContainsNMkfileActions(fops[0], 1)
				req.Equal("/dest/foo", mkfiles[0].Mkfile.Path)
				req.Equal([]byte(`some content`), mkfiles[0].Mkfile.Data)
				req.Equal(int32(0o0644), mkfiles[0].Mkfile.Mode)
			},
			testcompile.WithInitialState(llb.Scratch().Dir("/dest")),
		)

		compile.Test(
			"options/mode",
			`state.#File & { file: { mkfile: "./foo", content: 'some content', options: mode: 0o0666 } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkfiles := req.ContainsNMkfileActions(fops[0], 1)
				req.Equal(int32(0o0666), mkfiles[0].Mkfile.Mode)
			},
		)

		compile.Test(
			"options/ctime",
			`state.#File & { file: { mkfile: "./foo", content: 'some content', options: ctime: "2016-04-11T17:23:07-07:00" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkfiles := req.ContainsNMkfileActions(fops[0], 1)
				req.Equal(int64(1460420587000000000), mkfiles[0].Mkfile.Timestamp)
			},
		)

		compile.Test(
			"options/user/name",
			`state.#File & { file: { mkfile: "./foo", content: 'some content', options: user: "kim" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkfiles := req.ContainsNMkfileActions(fops[0], 1)
				req.Equal("kim", mkfiles[0].Mkfile.Owner.User.User.(*pb.UserOpt_ByName).ByName.Name)
			},
		)

		compile.Test(
			"options/user/uid",
			`state.#File & { file: { mkfile: "./foo", content: 'some content', options: uid: 9 } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkfiles := req.ContainsNMkfileActions(fops[0], 1)
				req.Equal(uint32(9), mkfiles[0].Mkfile.Owner.User.User.(*pb.UserOpt_ByID).ByID)
			},
		)

		compile.Test(
			"options/group/name",
			`state.#File & { file: { mkfile: "./foo", content: 'some content', options: group: "breeders" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkfiles := req.ContainsNMkfileActions(fops[0], 1)
				req.Equal("breeders", mkfiles[0].Mkfile.Owner.Group.User.(*pb.UserOpt_ByName).ByName.Name)
			},
		)

		compile.Test(
			"options/group/gid",
			`state.#File & { file: { mkfile: "./foo", content: 'some content', options: gid: 1 } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkfiles := req.ContainsNMkfileActions(fops[0], 1)
				req.Equal(uint32(1), mkfiles[0].Mkfile.Owner.Group.User.(*pb.UserOpt_ByID).ByID)
			},
		)
	})

	compile.Run("Mkdir", func(compile *testcompile.Tester) {
		compile.Test(
			"minimal",
			`state.#File & { file: { mkdir: "./foo" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkdirs := req.ContainsNMkdirActions(fops[0], 1)
				req.Equal("/dest/foo", mkdirs[0].Mkdir.Path)
				req.Equal(int32(0o0755), mkdirs[0].Mkdir.Mode)
				req.False(mkdirs[0].Mkdir.MakeParents)
			},
			testcompile.WithInitialState(llb.Scratch().Dir("/dest")),
		)

		compile.Test(
			"options/createParents",
			`state.#File & { file: { mkdir: "./foo", options: createParents: true } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkdirs := req.ContainsNMkdirActions(fops[0], 1)
				req.True(mkdirs[0].Mkdir.MakeParents)
			},
		)

		compile.Test(
			"options/mode",
			`state.#File & { file: { mkdir: "./foo", options: mode: 0o0777 } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkdirs := req.ContainsNMkdirActions(fops[0], 1)
				req.Equal(int32(0o0777), mkdirs[0].Mkdir.Mode)
			},
		)

		compile.Test(
			"options/ctime",
			`state.#File & { file: { mkdir: "./foo", options: ctime: "2016-04-11T17:23:07-07:00" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkdirs := req.ContainsNMkdirActions(fops[0], 1)
				req.Equal(int64(1460420587000000000), mkdirs[0].Mkdir.Timestamp)
			},
		)

		compile.Test(
			"options/user/name",
			`state.#File & { file: { mkdir: "./foo", options: user: "kim" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkdirs := req.ContainsNMkdirActions(fops[0], 1)
				req.Equal("kim", mkdirs[0].Mkdir.Owner.User.User.(*pb.UserOpt_ByName).ByName.Name)
			},
		)

		compile.Test(
			"options/user/uid",
			`state.#File & { file: { mkdir: "./foo", options: uid: 9 } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkdirs := req.ContainsNMkdirActions(fops[0], 1)
				req.Equal(uint32(9), mkdirs[0].Mkdir.Owner.User.User.(*pb.UserOpt_ByID).ByID)
			},
		)

		compile.Test(
			"options/group/name",
			`state.#File & { file: { mkdir: "./foo", options: group: "breeders" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkdirs := req.ContainsNMkdirActions(fops[0], 1)
				req.Equal("breeders", mkdirs[0].Mkdir.Owner.Group.User.(*pb.UserOpt_ByName).ByName.Name)
			},
		)

		compile.Test(
			"options/group/gid",
			`state.#File & { file: { mkdir: "./foo", options: gid: 1 } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, mkdirs := req.ContainsNMkdirActions(fops[0], 1)
				req.Equal(uint32(1), mkdirs[0].Mkdir.Owner.Group.User.(*pb.UserOpt_ByID).ByID)
			},
		)
	})

	compile.Run("Rm", func(compile *testcompile.Tester) {
		compile.Test(
			"minimal",
			`state.#File & { file: { rm: "./foo" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, rms := req.ContainsNRmActions(fops[0], 1)
				req.Equal("/dest/foo", rms[0].Rm.Path)
				req.False(rms[0].Rm.AllowNotFound)
				req.False(rms[0].Rm.AllowWildcard)
			},
			testcompile.WithInitialState(llb.Scratch().Dir("/dest")),
		)

		compile.Test(
			"options/allowNotFound",
			`state.#File & { file: { rm: "./foo", options: allowNotFound: true } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, rms := req.ContainsNRmActions(fops[0], 1)
				req.True(rms[0].Rm.AllowNotFound)
			},
		)

		compile.Test(
			"options/wildcard",
			`state.#File & { file: { rm: "./foo/*", options: wildcard: true } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, rms := req.ContainsNRmActions(fops[0], 1)
				req.True(rms[0].Rm.AllowWildcard)
			},
		)
	})
}
