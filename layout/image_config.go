package layout

import "gitlab.wikimedia.org/dduvall/phyton/common"

type ImageConfig struct {
	From             string
	User             common.User
	Environment      common.Env
	Entrypoint       []string
	DefaultArguments []string
	WorkingDirectory string
	Labels           map[string]string
	StopSignal       string
}
