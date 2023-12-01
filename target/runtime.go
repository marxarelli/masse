package target

import (
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type Runtime struct {
	User       string     `json:"user"`
	Env        common.Env `json:"env"`
	Entrypoint []string   `json:"entrypoint"`
	Arguments  []string   `json:"arguments"`
	Directory  string     `json:"directory"`
	StopSignal string     `json:"stopSignal"`
}

func (rt Runtime) ImageConfig() oci.ImageConfig {
	return oci.ImageConfig{
		User:       rt.User,
		Env:        rt.Env.Assignments(),
		Entrypoint: rt.Entrypoint,
		Cmd:        rt.Arguments,
		WorkingDir: rt.Directory,
		StopSignal: rt.StopSignal,
	}
}
