package state

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/moby/buildkit/client/llb"
)

type LLBRunOptions interface {
	LLBRunOptions(ChainStates) ([]llb.RunOption, error)
}

type Run struct {
	Command   string     `json:"run"`
	Arguments []string   `json:"arguments"`
	Options   RunOptions `json:"optionsValue"`
}

func (run *Run) Description() string {
	return fmt.Sprintf(
		"$ %s",
		run.ShellCommand(),
	)
}

func (run *Run) Compile(primary llb.State, secondary ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	ro, err := run.LLBRunOptions(secondary)
	if err != nil {
		return primary, err
	}

	return primary.Run(append(ro, constraintsTo[llb.RunOption](constraints)...)...).Root(), nil
}

func (run *Run) LLBRunOptions(states ChainStates) ([]llb.RunOption, error) {
	ro := []llb.RunOption{}
	ro = append(ro, llb.Args(run.ShellArgs()))

	oro, err := run.Options.LLBRunOptions(states)
	if err != nil {
		return ro, err
	}

	return append(ro, oro...), nil
}

func (run *Run) ChainRefs() []ChainRef {
	refs := []ChainRef{}

	for _, op := range run.Options {
		if op.SourceMount != nil {
			refs = append(refs, op.SourceMount.From)
		}
	}

	return refs
}

func (run *Run) ShellCommand() string {
	return run.Command + " " + strings.Join(run.QuotedArguments(), " ")
}

func (run *Run) ShellArgs() []string {
	return []string{"/bin/sh", "-c", run.ShellCommand()}
}

func (run *Run) QuotedArguments() []string {
	quoted := make([]string, len(run.Arguments))

	for i, arg := range run.Arguments {
		quoted[i] = strconv.Quote(arg)
	}

	return quoted
}

type RunOptions []*RunOption

func (opts RunOptions) LLBRunOptions(states ChainStates) ([]llb.RunOption, error) {
	ro := []llb.RunOption{}

	for _, opt := range opts {
		oro, err := opt.LLBRunOptions(states)
		if err != nil {
			return ro, err
		}

		ro = append(ro, oro...)
	}

	return ro, nil
}

type RunOption struct {
	*Host
	*CacheMount
	*SourceMount
	*TmpFSMount
	*ReadOnly
	*Option
}

func (opt *RunOption) LLBRunOptions(states ChainStates) ([]llb.RunOption, error) {
	llbOpt, ok := oneof[LLBRunOptions](opt)
	if ok {
		return llbOpt.LLBRunOptions(states)
	}

	return nil, nil
}
