package target

import "gitlab.wikimedia.org/dduvall/phyton/common"

type Runtime struct {
	User       common.User `json:"user"`
	Env        common.Env  `json:"env"`
	Entrypoint []string    `json:"entrypoint"`
	Arguments  []string    `json:"arguments"`
	Directory  string      `json:"directory"`
	StopSignal string      `json:"stopSignal"`
}
