package v1

import (
	"context"
	"testing"

	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
	"gitlab.wikimedia.org/dduvall/masse/util/testmetaresolver"
)

func TestImage(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/dduvall/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			return newCompiler(nil, WithImageMetaResolver(testmetaresolver.New(
				"an.example/image/ref",
				oci.Image{
					Config: oci.ImageConfig{
						WorkingDir: "/workdir",
					},
				},
			)))
		}),
	)

	compile.Test(
		"minimal",
		`state.#Image & { image: "an.example/image/ref" }`,
		func(t *testing.T, req *llbtest.Assertions, run *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("docker-image://an.example/image/ref:latest", sops[0].Source.Identifier)

			dir, err := run.State.GetDir(context.TODO())
			req.NoError(err)

			req.Equal("/workdir", dir)
		},
	)

	compile.Test(
		"inherit/false",
		`state.#Image & { image: "an.example/image/ref", inherit: false }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("docker-image://an.example/image/ref:latest", sops[0].Source.Identifier)
		},
	)

	compile.Test(
		"options/layerLimit",
		`state.#Image & { image: "an.example/image/ref", options: layerLimit: 20 }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("docker-image://an.example/image/ref:latest", sops[0].Source.Identifier)
			req.Contains(sops[0].Source.Attrs, "image.layerlimit")
			req.Equal("20", sops[0].Source.Attrs["image.layerlimit"])
		},
	)
}
