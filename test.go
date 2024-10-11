package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/tools/flow"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/imagemetaresolver"
	"github.com/pkg/errors"
)

type Chain []*State

type State struct {
	*Scratch `json:",inline"`
	*Image   `json:",inline"`
	*Run     `json:",inline"`
}

func (state *State) UnmarshalJSON(data []byte) error {
	st := map[string]json.RawMessage{}
	err := json.Unmarshal(data, &st)
	if err != nil {
		return err
	}

	if _, ok := st["scratch"]; ok {
		state.Scratch = &Scratch{Scratch: true}
		return nil
	}

	if _, ok := st["image"]; ok {
		state.Image = &Image{}
		return json.Unmarshal(data, state.Image)
	}

	if _, ok := st["run"]; ok {
		state.Run = &Run{}
		return json.Unmarshal(data, state.Run)
	}

	return nil
}

type Scratch struct {
	Scratch bool `json:"scratch"`
}

type Image struct {
	Ref     string `json:"image"`
	Inherit bool   `json:"inherit"`
}

type Run struct {
	Command   string     `json:"run"`
	Arguments []string   `json:"arguments"`
	Options   RunOptions `json:"optionsValue"`
}

type RunOptions []*RunOption

type RunOption struct {
	*CacheMount
	*SourceMount
	*TmpFSMount
	*ReadOnly
}

const (
	CacheShared  CacheAccess = "shared"
	CachePrivate             = "private"
	CacheLocked              = "locked"
)

type CacheAccess string

type CacheMount struct {
	Target string `json:"cache"`
	Access CacheAccess
}

type SourceMount struct {
	Target string `json:"mount"`
	From   Chain
	Source string
}

type TmpFSMount struct {
	TmpFS string
	Size  uint64
}

type ReadOnly struct {
	ReadOnly bool
}

func main() {
	path := "test.cue"

	ctx := cuecontext.New(cuecontext.EvaluatorVersion(cuecontext.EvalV2))

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	root := ctx.CompileBytes(data, cue.Filename(path))
	if err := root.Err(); err != nil {
		log.Fatal(err)
	}

	chains := root.LookupPath(cue.ParsePath("chains"))
	if chains.Err() != nil {
		log.Fatal(chains.Err())
	}

	controller := flow.New(nil, chains, func(v cue.Value) (flow.Runner, error) {
		path := v.Path()
		fmt.Printf("path: %s\n", path)

		if path.String() == "chains" {
			return nil, nil
		}

		return flow.RunnerFunc(func(t *flow.Task) error {
			fmt.Printf("task: %#v\n", t.Path())
			fmt.Printf("dependencies: %#v\n", t.Dependencies())

			_, err := compileChain(t.Value())

			return err
		}), nil
	})

	if err := controller.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func compileChain(v cue.Value) (llb.State, error) {
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
		state, err = compileState(state, list.Value())

		if err != nil {
			return state, err
		}
	}

	return state, err
}

type compiler func(llb.State, cue.Value) (llb.State, error)

func compileState(state llb.State, v cue.Value) (llb.State, error) {
	var stateCompilers = map[string]compiler{
		"scratch": compileScratch,
		"image":   compileImage,
		"run":     compileRun,
	}

	for key, compiler := range stateCompilers {
		if v.LookupPath(cue.ParsePath(key)).Exists() {
			return compiler(state, v)
		}
	}

	return state, errorf(v, "unsupported state operation")
}

func compileScratch(state llb.State, v cue.Value) (llb.State, error) {
	return llb.Scratch(), nil
}

func compileImage(state llb.State, v cue.Value) (llb.State, error) {
	ref := requireString(v, "image")
	inherit := requireBool(v, "inherit")

	if inherit {
		return llb.Image(ref, llb.WithMetaResolver(imagemetaresolver.Default())), nil
	}

	return llb.Image(ref), nil
}

type runOptionCompiler func(cue.Value) (llb.RunOption, error)

func compileRun(state llb.State, v cue.Value) (llb.State, error) {
	var runOptionCompilers = map[string]runOptionCompiler{
		"cache":    compileCache,
		"mount":    compileMount,
		"tmpfs":    compileTmpfs,
		"readonly": compileReadonly,
	}

	cmd := requireString(v, "run")

	args := []string{}
	argv := lookup(v, "arguments")

	if argv.Exists() {
		list, err := argv.List()
		if err != nil {
			return state, vError(v, err)
		}

		i := 0
		for list.Next() {
			arg, err := list.Value().String()
			if err != nil {
				return state, vError(v, err)
			}

			args[i] = strconv.Quote(arg)
			i++
		}
	}

	optsv := lookup(v, "optionsValue")
	options := []llb.RunOption{
		llb.Args([]string{
			"/bin/sh",
			"-c",
			cmd + " " + strings.Join(args, " "),
		}),
	}

	if optsv.Exists() {
		list, err := optsv.List()
		if err != nil {
			return state, vError(v, err)
		}

		for list.Next() {
			ov := list.Value()
			for key, compiler := range runOptionCompilers {
				if ov.LookupPath(cue.ParsePath(key)).Exists() {
					ro, err := compiler(ov)
					if err != nil {
						return state, vError(v, err)
					}

					if ro != nil {
						options = append(options, ro)
					}
					break
				}
			}
		}
	}

	return state.Run(options...).Root(), nil
}

func compileCache(v cue.Value) (llb.RunOption, error) {
	cache := requireString(v, "cache")
	access := requireString(v, "access")

	return llb.AddMount(
		cache,
		llb.Scratch(),
		llb.AsPersistentCacheDir(cache, compileCacheAccess(access)),
	), nil
}

func compileCacheAccess(access string) llb.CacheMountSharingMode {
	switch access {
	case "private":
		return llb.CacheMountPrivate
	case "locked":
		return llb.CacheMountLocked
	}

	return llb.CacheMountShared
}

func compileMount(v cue.Value) (llb.RunOption, error) {
	target := requireString(v, "mount")
	source := requireString(v, "source")
	state := llb.Scratch()
	from := v.LookupPath(cue.ParsePath("from"))
	if from.Exists() {
		var err error
		state, err = compileChain(from)
		if err != nil {
			return nil, err
		}
	}

	return llb.AddMount(target, state, llb.SourcePath(source)), nil
}

func compileTmpfs(v cue.Value) (llb.RunOption, error) {
	target := requireString(v, "tmpfs")
	size := requireInt64(v, "size")

	return llb.AddMount(
		target,
		llb.Scratch(),
		llb.Tmpfs(llb.TmpfsSize(size)),
	), nil
}

func compileReadonly(v cue.Value) (llb.RunOption, error) {
	readonly := requireBool(v, "readonly")

	if readonly {
		return llb.ReadonlyRootFS(), nil
	}

	return nil, nil
}

func lookup(root cue.Value, path string) cue.Value {
	return root.LookupPath(cue.ParsePath(path))
}

func requireBool(root cue.Value, path string) bool {
	v := lookup(root, path)
	b, err := v.Bool()
	if err != nil {
		log.Panic(vError(v, err))
	}

	return b
}

func requireString(root cue.Value, path string) string {
	v := lookup(root, path)
	str, err := v.String()
	if err != nil {
		log.Panic(vError(v, err))
	}

	return str
}

func requireInt64(root cue.Value, path string) int64 {
	v := lookup(root, path)
	i, err := v.Int64()
	if err != nil {
		log.Panic(vError(v, err))
	}

	return i
}

func vError(v cue.Value, err error) error {
	return errorf(v, "compile error: %s", err)
}

func errorf(v cue.Value, msg string, args ...any) error {
	return errors.Errorf(fmt.Sprintf("%s: %s at %s", v.Path(), msg, v.Pos()), args...)
}
