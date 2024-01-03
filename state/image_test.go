package state

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/masse/common"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testdecode"
)

func TestDecodeImage(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/masse/schema/state"},
	}

	testdecode.Run(tester,
		"state.#Image",
		`state.#Image & { image: "foo.example/image/ref" }`,
		Image{
			Ref:     "foo.example/image/ref",
			Inherit: true,
		},
	)

	testdecode.Run(tester,
		"state.#Image/inherit:false",
		`state.#Image & { image: "foo.example/image/ref", inherit: false }`,
		Image{
			Ref:     "foo.example/image/ref",
			Inherit: false,
		},
	)

	testdecode.Run(tester,
		"state.#Image/options/platform/literal",
		`state.#Image & { image: "foo.example/image/ref", options: [ { platform: { os: "linux", architecture: "amd64" } } ] }`,
		Image{
			Ref:     "foo.example/image/ref",
			Inherit: true,
			Options: []*ImageOption{
				{Constraint: &Constraint{Platform: &Platform{Platform: common.Platform{
					OS:           "linux",
					Architecture: "amd64",
				}}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Image/options/platform/symbolic",
		`state.#Image & { image: "foo.example/image/ref", options: [ { platform: "linux/amd64" } ] }`,
		Image{
			Ref:     "foo.example/image/ref",
			Inherit: true,
			Options: []*ImageOption{
				{Constraint: &Constraint{Platform: &Platform{Platform: common.Platform{
					OS:           "linux",
					Architecture: "amd64",
				}}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Image/options/layerLimit",
		`state.#Image & { image: "foo.example/image/ref", options: [ { layerLimit: 999 } ] }`,
		Image{
			Ref:     "foo.example/image/ref",
			Inherit: true,
			Options: []*ImageOption{
				{LayerLimit: &LayerLimit{LayerLimit: 999}},
			},
		},
	)
}

func TestCompileImage(t *testing.T) {
	req := require.New(t)

	image := Image{
		Ref: "some.example/image:ref",
		Options: ImageOptions{
			{LayerLimit: &LayerLimit{LayerLimit: 999}},
			{Constraint: &Constraint{Platform: &Platform{Platform: common.Platform{
				OS:           "linux",
				Architecture: "arm64",
				Variant:      "v8",
			}}}},
		},
	}

	state, err := image.CompileSource(ChainStates{})
	req.NoError(err)

	def, err := state.Marshal(context.TODO())
	req.NoError(err)

	llbreq := llbtest.New(t, def)

	ops, sops := llbreq.ContainsNSourceOps(1)
	req.Equal("docker-image://some.example/image:ref", sops[0].Source.Identifier)

	req.Contains(sops[0].Source.Attrs, "image.layerlimit")
	req.Equal("999", sops[0].Source.Attrs["image.layerlimit"])

	platform := ops[0].GetPlatform()
	req.NotNil(platform)
	req.Equal("linux", platform.OS)
	req.Equal("arm64", platform.Architecture)
	req.Equal("v8", platform.Variant)
}
