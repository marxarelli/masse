package state

import (
	"context"
	"testing"

	"github.com/moby/buildkit/client/llb"
	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/phyton/common"
	"gitlab.wikimedia.org/dduvall/phyton/util/llbtest"
	"gitlab.wikimedia.org/dduvall/phyton/util/testdecode"
)

func TestDecodeLocal(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/phyton/schema/state"},
	}

	testdecode.Run(tester,
		"state.#Local",
		`state.#Local & { local: "context" }`,
		Local{Name: "context"},
	)

	testdecode.Run(tester,
		"state.#Local/options/differ",
		`state.#Local & { local: "context", options: [ { differ: "metadata", require: true } ] }`,
		Local{
			Name: "context",
			Options: []*LocalOption{
				{Differ: &Differ{Differ: "metadata", Require: true}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Local/options/exclude",
		`state.#Local & { local: "context", options: [ { exclude: ["/src/foo"] } ] }`,
		Local{
			Name: "context",
			Options: []*LocalOption{
				{Exclude: &Exclude{Exclude: []common.Glob{"/src/foo"}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Local/options/followPaths",
		`state.#Local & { local: "context", options: [ { followPaths: ["/src/foo"] } ] }`,
		Local{
			Name: "context",
			Options: []*LocalOption{
				{FollowPaths: &FollowPaths{FollowPaths: []string{"/src/foo"}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Local/options/include",
		`state.#Local & { local: "context", options: [ { include: ["/src/foo"] } ] }`,
		Local{
			Name: "context",
			Options: []*LocalOption{
				{Include: &Include{Include: []common.Glob{"/src/foo"}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Local/options/sharedKeyHint",
		`state.#Local & { local: "context", options: [ { sharedKeyHint: "foo-hint" } ] }`,
		Local{
			Name: "context",
			Options: []*LocalOption{
				{SharedKeyHint: &SharedKeyHint{SharedKeyHint: "foo-hint"}},
			},
		},
	)
}

func TestCompileLocal(t *testing.T) {
	req := require.New(t)

	local := Local{
		Name: "context",
		Options: LocalOptions{
			{Include: &Include{Include: []common.Glob{"foo/*", "*.bar"}}},
			{Exclude: &Exclude{Exclude: []common.Glob{"baz*"}}},
			{FollowPaths: &FollowPaths{FollowPaths: []string{"foo/path"}}},
			{SharedKeyHint: &SharedKeyHint{SharedKeyHint: "foo-hint"}},
			{Differ: &Differ{Differ: DiffNone, Require: true}},
		},
	}

	state, err := local.Compile(llb.Scratch(), ChainStates{})
	req.NoError(err)

	def, err := state.Marshal(context.TODO())
	req.NoError(err)

	llbreq := llbtest.New(t, def)

	_, sops := llbreq.ContainsNSourceOps(1)
	req.Equal("local://context", sops[0].Source.Identifier)

	req.Contains(sops[0].Source.Attrs, "local.includepattern")
	req.Equal(`["foo/*","*.bar"]`, sops[0].Source.Attrs["local.includepattern"])

	req.Contains(sops[0].Source.Attrs, "local.excludepatterns")
	req.Equal(`["baz*"]`, sops[0].Source.Attrs["local.excludepatterns"])

	req.Contains(sops[0].Source.Attrs, "local.followpaths")
	req.Equal(`["foo/path"]`, sops[0].Source.Attrs["local.followpaths"])

	req.Contains(sops[0].Source.Attrs, "local.sharedkeyhint")
	req.Equal("foo-hint", sops[0].Source.Attrs["local.sharedkeyhint"])

	req.Contains(sops[0].Source.Attrs, "local.differ")
	req.Equal("none", sops[0].Source.Attrs["local.differ"])
}
