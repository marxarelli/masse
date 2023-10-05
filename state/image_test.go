package state

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/phyton/common"
	"gitlab.wikimedia.org/dduvall/phyton/util/testdecode"
)

func TestDecodeImage(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/phyton/schema/state"},
	}

	testdecode.Run(tester,
		"state.#Image",
		`state.#Image & { image: "foo.example/image/ref" }`,
		Image{
			Ref: "foo.example/image/ref",
		},
	)

	testdecode.Run(tester,
		"state.#Image/options/platform/literal",
		`state.#Image & { image: "foo.example/image/ref", options: [ { platform: { os: "linux", architecture: "amd64" } } ] }`,
		Image{
			Ref: "foo.example/image/ref",
			Options: []*ImageOption{
				{Platform: &Platform{Platform: common.Platform{OS: "linux", Architecture: "amd64"}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Image/options/platform/symbolic",
		`state.#Image & { image: "foo.example/image/ref", options: [ { platform: "linux/amd64" } ] }`,
		Image{
			Ref: "foo.example/image/ref",
			Options: []*ImageOption{
				{Platform: &Platform{Platform: common.Platform{OS: "linux", Architecture: "amd64"}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Image/options/layerLimit",
		`state.#Image & { image: "foo.example/image/ref", options: [ { layerLimit: 999 } ] }`,
		Image{
			Ref: "foo.example/image/ref",
			Options: []*ImageOption{
				{LayerLimit: &LayerLimit{LayerLimit: 999}},
			},
		},
	)
}
