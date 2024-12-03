package v1

import (
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileRun(state llb.State, v cue.Value) (llb.State, error) {
	cmd, err := lookup.String(v, "run")
	if err != nil {
		return state, vError(v, err)
	}

	args, err := lookup.DecodeListOrSingle[[]string](v, "arguments")
	if err != nil {
		return state, vError(v, err)
	}

	for i, arg := range args {
		args[i] = strconv.Quote(arg)
	}

	options, err := lookup.DecodeOptions[RunOptions](v)
	if err != nil {
		return state, vError(v, err)
	}

	options.SetCompiler(c)

	if len(args) > 0 {
		cmd += " " + strings.Join(args, " ")
	}

	llbOptions := []llb.RunOption{
		llb.Args([]string{"/bin/sh", "-c", cmd}),
		options,
	}

	return state.Run(llbOptions...).Root(), nil
}

type RunOptions []*RunOption

func (opts RunOptions) SetRunOption(info *llb.ExecInfo) {
	for _, opt := range opts {
		opt.SetRunOption(info)
	}
}

func (opts RunOptions) SetCompiler(c *compiler) {
	for _, opt := range opts {
		opt.SetCompiler(c)
	}
}

type RunOption struct {
	*Host
	*CacheMount
	*SourceMount
	*TmpFSMount
	*ValidExitCodes
	*Option
	*Constraint
}

func (opt *RunOption) SetRunOption(info *llb.ExecInfo) {
	withOneOf(opt, func(ro llb.RunOption) { ro.SetRunOption(info) })
}

func (opt *RunOption) SetCompiler(c *compiler) {
	withOneOf(opt, func(subc subcompiler) { subc.SetCompiler(c) })
}

type Host common.Host

func (host *Host) SetRunOption(info *llb.ExecInfo) {
	llb.AddExtraHost(host.Host, host.IP).SetRunOption(info)
}

const (
	CacheShared  CacheAccess = "shared"
	CachePrivate             = "private"
	CacheLocked              = "locked"
)

type CacheAccess string

type CacheMount struct {
	Target string      `json:"cache"`
	Access CacheAccess `json:"access"`
}

func (cm *CacheMount) SetRunOption(info *llb.ExecInfo) {
	mode := llb.CacheMountShared

	switch cm.Access {
	case CachePrivate:
		mode = llb.CacheMountPrivate
	case CacheLocked:
		mode = llb.CacheMountLocked
	}

	llb.AddMount(
		cm.Target,
		llb.Scratch(),
		llb.AsPersistentCacheDir(cm.Target, mode),
	).SetRunOption(info)
}

type SourceMount struct {
	Target   string    `json:"mount"`
	From     cue.Value `json:"from"`
	Source   string    `json:"source"`
	Readonly bool      `json:"readonly"`
	c        *compiler
}

func (sm *SourceMount) SetCompiler(c *compiler) {
	sm.c = c
}

func (sm *SourceMount) SetRunOption(info *llb.ExecInfo) {
	from, _ := sm.c.compileChainByRef(sm.From)

	opts := []llb.MountOption{
		llb.SourcePath(sm.Source),
	}

	if sm.Readonly {
		opts = append(opts, llb.Readonly)
	}

	llb.AddMount(
		sm.Target,
		from,
		opts...,
	).SetRunOption(info)
}

type TmpFSMount struct {
	TmpFS string `json:"tmpfs"`
	Size  uint64 `json:"size"`
}

func (tm *TmpFSMount) SetRunOption(info *llb.ExecInfo) {
	llb.AddMount(
		tm.TmpFS,
		llb.Scratch(),
		llb.Tmpfs(llb.TmpfsSize(int64(tm.Size))),
	).SetRunOption(info)
}

type ValidExitCodes struct {
	ValidExitCodes []int `json:"validExitCodes"`
}

func (vec *ValidExitCodes) SetRunOption(info *llb.ExecInfo) {
	llb.ValidExitCodes(vec.ValidExitCodes...).SetRunOption(info)
}
