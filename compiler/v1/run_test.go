package v1

import (
	"testing"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestRun(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/dduvall/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			return newCompiler(nil)
		}),
	)

	compile.Test(
		"minimal",
		`state.#Run & { run: "make" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Equal([]string{"/bin/sh", "-c", "make"}, eops[0].Exec.Meta.Args)
		},
	)

	compile.Test(
		"arguments/single",
		`state.#Run & { run: "make", arguments: "foo bar" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Equal([]string{"/bin/sh", "-c", `make "foo bar"`}, eops[0].Exec.Meta.Args)
		},
	)

	compile.Test(
		"arguments/multiple",
		`state.#Run & { run: "make", arguments: ["foo", "bar"] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Equal([]string{"/bin/sh", "-c", `make "foo" "bar"`}, eops[0].Exec.Meta.Args)
		},
	)

	compile.Test(
		"defaultOptions/customName",
		`state.#Run & { run: "make" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, _ := req.ContainsNExecOps(1)
			md := req.HasMetadata(ops[0])
			req.Contains(md.Description, "llb.customname")
			req.Equal("💻 make", md.Description["llb.customname"])
		},
	)

	compile.Test(
		"defaultOptions/customName/arguments",
		`state.#Run & { run: "make", arguments: ["foo"] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, _ := req.ContainsNExecOps(1)
			md := req.HasMetadata(ops[0])
			req.Contains(md.Description, "llb.customname")
			req.Equal("💻 make foo", md.Description["llb.customname"])
		},
	)

	compile.Test(
		"options/host",
		`state.#Run & { run: "make", options: [ { host: "foo", ip: "1.1.1.1" } ] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.NotNil(eops[0].Exec.Meta.ExtraHosts)
			req.Len(eops[0].Exec.Meta.ExtraHosts, 1)
			req.Equal("foo", eops[0].Exec.Meta.ExtraHosts[0].Host)
			req.Equal("1.1.1.1", eops[0].Exec.Meta.ExtraHosts[0].IP)
		},
	)

	compile.Test(
		"options/cache",
		`state.#Run & { run: "apt-get install", options: [ { cache: "/var/lib/apt" } ] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Len(eops[0].Exec.Mounts, 2)
			mnt := eops[0].Exec.Mounts[1]

			req.Equal("/var/lib/apt", mnt.Dest)
			req.Equal(pb.MountType_CACHE, mnt.MountType)
			req.NotNil(mnt.CacheOpt)
			req.Equal(pb.CacheSharingOpt_SHARED, mnt.CacheOpt.Sharing)
		},
	)

	compile.Test(
		"options/mount",
		`state.#Run & { run: "make -C /src", options: [ { mount: "/src", from: "repo", source: "/component/src" } ] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, eops := req.ContainsNExecOps(1)
			req.Len(eops[0].Exec.Mounts, 2)
			mnt := eops[0].Exec.Mounts[1]

			req.Equal("/src", mnt.Dest)
			req.Equal(pb.MountType_BIND, mnt.MountType)
			req.Equal("/component/src", mnt.Selector)
			req.False(mnt.Readonly)

			_, sops := req.ContainsNSourceOps(1)
			req.Equal("git://an.example/repo.git#refs/heads/main", sops[0].Source.Identifier)

			iops := req.HasValidInputs(ops[0])
			req.Equal(iops[0].Op, sops[0])
		},
		testcompile.WithCompiler(func() *compiler {
			c := newCompiler(nil)
			c.chainCompilers = map[string]chainCompiler{
				"repo": func(_ *compiler) *chainResult {
					return &chainResult{state: llb.Git("an.example/repo.git", "refs/heads/main")}
				},
			}
			return c
		}),
	)

	compile.Test(
		"options/mount/readonly",
		`state.#Run & { run: "scan /src", options: [ { mount: "/src", from: "repo", readonly: true } ] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, eops := req.ContainsNExecOps(1)
			req.Len(eops[0].Exec.Mounts, 2)
			mnt := eops[0].Exec.Mounts[1]

			req.Equal("/src", mnt.Dest)
			req.Equal(pb.MountType_BIND, mnt.MountType)
			req.Equal("/", mnt.Selector)
			req.True(mnt.Readonly)

			_, sops := req.ContainsNSourceOps(1)
			req.Equal("git://an.example/repo.git#refs/heads/main", sops[0].Source.Identifier)

			iops := req.HasValidInputs(ops[0])
			req.Equal(iops[0].Op, sops[0])
		},
		testcompile.WithCompiler(func() *compiler {
			c := newCompiler(nil)
			c.chainCompilers = map[string]chainCompiler{
				"repo": func(_ *compiler) *chainResult {
					return &chainResult{state: llb.Git("an.example/repo.git", "refs/heads/main")}
				},
			}
			return c
		}),
	)

	compile.Test(
		"options/tmpfs",
		`state.#Run & { run: "make -C /src", options: [ { tmpfs: "/tmp" } ] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Len(eops[0].Exec.Mounts, 2)
			mnt := eops[0].Exec.Mounts[1]

			req.Equal("/tmp", mnt.Dest)
			req.Equal(pb.MountType_TMPFS, mnt.MountType)
			req.NotNil(mnt.TmpfsOpt)
			req.Equal(int64(1024*1024*100), mnt.TmpfsOpt.Size)
		},
	)

	compile.Test(
		"options/tmpfs/size",
		`state.#Run & { run: "make -C /src", options: [ { tmpfs: "/tmp", size: 200Mi } ] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Len(eops[0].Exec.Mounts, 2)
			mnt := eops[0].Exec.Mounts[1]

			req.Equal("/tmp", mnt.Dest)
			req.Equal(pb.MountType_TMPFS, mnt.MountType)
			req.NotNil(mnt.TmpfsOpt)
			req.Equal(int64(1024*1024*200), mnt.TmpfsOpt.Size)
		},
	)

	compile.Test(
		"options/customName",
		`state.#Run & { run: "make -C /src", options: customName: "building foo" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, _ := req.ContainsNExecOps(1)
			md := req.HasMetadata(ops[0])
			req.Contains(md.Description, "llb.customname")
			req.Equal("building foo", md.Description["llb.customname"])
		},
	)
}
