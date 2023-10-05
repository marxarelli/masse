package state

import "gitlab.wikimedia.org/dduvall/phyton/common"

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

type Env struct {
	Env common.Env
}

type Host common.Host

type Option struct {
	*Env
	*WorkingDirectory
}

type ReadOnly struct {
	ReadOnly bool
}

type TmpFSMount struct {
	TmpFS string
	Size  uint64
}

type WorkingDirectory struct {
	Directory string
}
