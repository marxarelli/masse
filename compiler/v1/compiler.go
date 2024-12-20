package v1

import (
	"context"
	errs "errors"
	"path/filepath"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/target"
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
		chainCache:     map[string]*chainResult{},
		chainCompilers: map[string]chainCompiler{},
		config:         cfg,
		errors:         []error{},
		ctx:            cfg.InitialContext,
	}
}

type compiler struct {
	chains         map[string]cue.Value
	chainCache     map[string]*chainResult
	chainCompilers map[string]chainCompiler
	config         *Config
	errors         []error
	mutex          sync.Mutex
	ctx            context.Context
	refStack       []string
}

func (c *compiler) copy(mut func(*compiler)) *compiler {
	newC := &compiler{
		chains:         c.chains,
		chainCache:     c.chainCache,
		chainCompilers: c.chainCompilers,
		config:         c.config,
		errors:         c.errors,
		ctx:            c.ctx,
		mutex:          c.mutex,
		refStack:       c.refStack,
	}
	mut(newC)
	return newC
}

type chainCompiler func(c *compiler) *chainResult

type chainResult struct {
	state llb.State
	err   error
}

func (c *compiler) Compile(target *target.Target) (llb.State, error) {
	c.defineChainCompilers()
	return c.compileChainByRef(target.Build)
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
				var chainMutex sync.Mutex
				chainMutex.Lock()
				defer chainMutex.Unlock()

				if result, cached := c.chainCache[ref]; cached {
					return result
				}

				st, err := c.compileChain(chain)
				result := &chainResult{state: st, err: err}
				c.chainCache[ref] = result
				return result
			}
		}(ref, chain)
	}
}

func (c *compiler) withRefOnStack(chainRef string) *compiler {
	return c.copy(func(c *compiler) {
		c.refStack = append(c.refStack, chainRef)
	})
}

func (c *compiler) constraints() Constraints {
	return c.config.Constraints()
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
