package v1

import (
	"context"
	errs "errors"
	"path/filepath"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/masse/target"
)

func New(chains map[string]cue.Value, options ...CompilerOption) target.Compiler {
	return newCompiler(chains, options...)
}

func newCompiler(chains map[string]cue.Value, options ...CompilerOption) *compiler {
	if chains == nil {
		chains = map[string]cue.Value{}
	}

	return &compiler{
		chains: chains,
		config: newConfig(options),
		errors: []error{},
		ctx:    context.Background(),
	}
}

type compiler struct {
	chains         map[string]cue.Value
	chainCompilers map[string]chainCompiler
	config         *Config
	errors         []error
	mutex          sync.Mutex
	ctx            context.Context
}

type chainCompiler func() *chainResult

type chainResult struct {
	state llb.State
	err   error
}

func (c *compiler) Compile(target *target.Target) (llb.State, error) {
	c.chainCompilers = map[string]chainCompiler{}

	for ref, chain := range c.chains {
		func(ref string, chain cue.Value) {
			c.chainCompilers[ref] = sync.OnceValue(func() *chainResult {
				st, err := c.compileChain(chain)
				return &chainResult{state: st, err: err}
			})
		}(ref, chain)
	}

	return c.compileChain(target.Build)
}

func (c *compiler) CompileState(state llb.State, v cue.Value) (llb.State, error) {
	return c.compileState(state, v)
}

func (c *compiler) Error() error {
	return errs.Join(c.errors...)
}

func (c *compiler) WithContext(ctx context.Context) target.Compiler {
	return c.withContext(ctx)
}

func (c *compiler) withContext(ctx context.Context) *compiler {
	return &compiler{
		chains:         c.chains,
		chainCompilers: c.chainCompilers,
		config:         c.config,
		errors:         c.errors,
		ctx:            ctx,
	}
}

func (c *compiler) compileChain(v cue.Value) (llb.State, error) {
	var err error
	state := llb.Scratch()

	err = v.Null()
	if err == nil {
		return state, nil
	}

	list, err := v.List()
	if err != nil {
		return state, vError(v, err)
	}

	for list.Next() {
		state, err = c.compileState(state, list.Value())

		if err != nil {
			return state, err
		}
	}

	return state, err
}

func (c *compiler) compileChainByRef(refv cue.Value) (llb.State, error) {
	state := llb.NewState(nil)

	ref, err := refv.String()
	if err != nil {
		return state, c.addVError(refv, err)
	}

	cc, ok := c.chainCompilers[ref]
	if !ok {
		return state, c.addError(errors.Errorf("unknown chain %q", ref))
	}

	result := cc()
	return result.state, c.addVError(refv, result.err)
}

func (c *compiler) absPath(state llb.State, path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	hadTrailingSlash := strings.HasSuffix(path, "/")

	cwd, _ := state.GetDir(c.ctx)
	abs := filepath.Join(cwd, path)

	if hadTrailingSlash {
		abs += "/"
	}

	return abs
}

func (c *compiler) addError(err error) error {
	if err != nil {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.errors = append(c.errors, err)
	}

	return err
}

func (c *compiler) addVError(v cue.Value, err error) error {
	return c.addError(vError(v, err))
}
