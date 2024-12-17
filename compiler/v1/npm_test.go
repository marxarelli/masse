package v1

import (
	"testing"

	"github.com/moby/buildkit/solver/pb"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestNpm(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/dduvall/masse/npm"},
		testcompile.WithCompiler(func() *compiler {
			return newCompiler(nil)
		}),
	)

	compile.Test(
		"default",
		`npm.install`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(2)
			req.Equal([]string{"/bin/sh", "-c", `npm install`}, eops[0].Exec.Meta.Args)
			req.Contains(eops[0].Exec.Meta.Env, "NPM_CONFIG_CACHE=/var/lib/cache/npm")
			req.Len(eops[0].Exec.Mounts, 2)
			mnt := eops[0].Exec.Mounts[1]
			req.Equal("/var/lib/cache/npm", mnt.Dest)
			req.Equal(pb.MountType_CACHE, mnt.MountType)
			req.NotNil(mnt.CacheOpt)
			req.Equal(pb.CacheSharingOpt_LOCKED, mnt.CacheOpt.Sharing)

			req.Equal([]string{"/bin/sh", "-c", `npm dedupe`}, eops[1].Exec.Meta.Args)
			req.Contains(eops[1].Exec.Meta.Env, "NPM_CONFIG_CACHE=/var/lib/cache/npm")
			req.Len(eops[1].Exec.Mounts, 2)
			mnt = eops[1].Exec.Mounts[1]
			req.Equal("/var/lib/cache/npm", mnt.Dest)
			req.Equal(pb.MountType_CACHE, mnt.MountType)
			req.NotNil(mnt.CacheOpt)
			req.Equal(pb.CacheSharingOpt_LOCKED, mnt.CacheOpt.Sharing)
		},
	)

	compile.Test(
		"only",
		`npm.install & { #only: "production" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(2)
			req.Equal([]string{"/bin/sh", "-c", `npm install "--only=production"`}, eops[0].Exec.Meta.Args)
		},
	)

	compile.Test(
		"options",
		`npm.install & { #options: directory: "/srv" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(2)
			req.Equal("/srv", eops[0].Exec.Meta.Cwd)
			req.Equal("/srv", eops[1].Exec.Meta.Cwd)
		},
	)
}
