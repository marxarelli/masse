package v1

import (
	"context"
	errs "errors"
	"path/filepath"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
	"gitlab.wikimedia.org/dduvall/masse/target"
	"golang.org/x/sync/singleflight"
)

func New(chains map[string]cue.Value, options ...CompilerOption) target.Compiler {
	return newCompiler(chains, options...)
}

func newCompiler(chains map[string]cue.Value, options ...CompilerOption) *compiler {
	if chains == nil {
		chains = map[string]cue.Value{}
	}

	cfg := newConfig(options)

	return &compiler{
		chains:         chains,
		chainCompilers: map[string]chainCompiler{},
		config:         cfg,
		errors:         []error{},
		ctx:            cfg.InitialContext,
		group:          new(singleflight.Group),
		states:         new(sync.Map),
	}
}

type compiler struct {
	chains         map[string]cue.Value
	chainCompilers map[string]chainCompiler
	config         *Config
	errors         []error
	ctx            context.Context
	refStack       []string
	group          *singleflight.Group
	states         *sync.Map
	mutex          sync.Mutex
}

func (c *compiler) copy(mut func(*compiler)) *compiler {
	newC := &compiler{
		chains:         c.chains,
		chainCompilers: c.chainCompilers,
		config:         c.config,
		errors:         c.errors,
		ctx:            c.ctx,
		refStack:       c.refStack,
		group:          c.group,
		mutex:          c.mutex,
		states:         c.states,
	}
	mut(newC)
	return newC
}

type chainCompiler func(c *compiler) *chainResult

type chainResult struct {
	state llb.State
	err   error
}

func (c *compiler) Compile(target *target.Target) (target.CompilerResult, error) {
	c.defineChainCompilers()

	state, err := c.compileChainByRef(target.Build)
	if err != nil {
		return nil, err
	}

	targetRef := lookup.NormalizeReference(target.Build)
	deps := map[string]llb.State{}
	c.states.Range(func(key, value any) bool {
		ref := key.(string)
		if ref != targetRef {
			deps[ref] = value.(llb.State)
		}
		return true
	})

	return &result{
		ref:      lookup.NormalizeReference(target.Build),
		platform: c.config.Platform.Platform,
		state:    state,
		deps:     deps,
	}, nil
}

func (c *compiler) CompileState(state llb.State, v cue.Value) (llb.State, error) {
	return c.compileState(state, v)
}

func (c *compiler) CompileChain(v cue.Value) (llb.State, error) {
	return c.compileChain(v)
}

func (c *compiler) Error() error {
	return errs.Join(c.errors...)
}

func (c *compiler) WithContext(ctx context.Context) target.Compiler {
	return c.copy(func(c *compiler) {
		c.ctx = ctx
	})
}

func (c *compiler) defineChainCompilers() {
	for ref, chain := range c.chains {
		func(ref string, chain cue.Value) {
			c.chainCompilers[ref] = func(c *compiler) *chainResult {
				st, err := c.compileChain(chain)
				return &chainResult{state: st, err: err}
			}
		}(ref, chain)
	}
}

func (c *compiler) withRefOnStack(chainRef string) *compiler {
	return c.copy(func(c *compiler) {
		c.refStack = append(c.refStack, chainRef)
	})
}

func (c *compiler) sourceConstraints() Constraints {
	return append(c.config.OpConstraints(), c.config.SourceConstraints()...)
}

func (c *compiler) opConstraints() Constraints {
	return c.config.OpConstraints()
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
