package testcompile

import (
	"context"
	"testing"

	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/masse/target"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testconfig"
)

type Tester struct {
	*testconfig.Tester
	options []TestOption
}

func New(t *testing.T, imports []string, options ...TestOption) *Tester {
	return &Tester{
		Tester:  testconfig.New(t, imports),
		options: options,
	}
}

type Test struct {
	Compiler target.Compiler
	State    llb.State
	Value    cue.Value
	Compile  CompileFunc
}

type CompileFunc func(*Test) (llb.State, error)

type TestOption func(*Test)

func WithInitialState(state llb.State) TestOption {
	return TestOption(func(test *Test) {
		test.State = state
	})
}

func WithCompiler[C target.Compiler](f func() C) TestOption {
	return TestOption(func(test *Test) {
		test.Compiler = f()
	})
}

func WithCompileFunc(f CompileFunc) TestOption {
	return TestOption(func(test *Test) {
		test.Compile = f
	})
}

func (tester *Tester) Run(name string, f func(*Tester)) {
	tester.Helper()

	tester.Tester.Run(name, func(t *testconfig.Tester) {
		f(&Tester{
			Tester:  t,
			options: tester.options,
		})
	})
}

func (tester *Tester) Test(
	name, expr string,
	testFunc func(*testing.T, *llbtest.Assertions, *Test),
	options ...TestOption,
) {
	tester.Helper()

	tester.Tester.Test(name, expr, func(t *testing.T, req *require.Assertions, v cue.Value) {
		t.Helper()

		test := &Test{
			State: llb.Scratch(),
			Value: v,
			Compile: func(test *Test) (llb.State, error) {
				return test.Compiler.CompileState(test.State, test.Value)
			},
		}

		for _, opt := range append(tester.options, options...) {
			opt(test)
		}

		if test.Compiler == nil {
			tester.Fatal(
				"compiler is nil. use WithCompiler() to setup a compiler when creating the tester or on an individual test",
			)
		}

		state, err := test.Compile(test)
		req.NoError(err)

		test.State = state

		def, err := state.Marshal(context.TODO())
		req.NoError(err)

		llbreq := llbtest.New(t, def)

		testFunc(t, llbreq, test)
	})
}
