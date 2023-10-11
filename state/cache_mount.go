package state

import "github.com/moby/buildkit/client/llb"

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

func (cm *CacheMount) LLBRunOptions(_ ChainStates) ([]llb.RunOption, error) {
	return []llb.RunOption{
		llb.AddMount(
			cm.Target,
			llb.Scratch(),
			llb.AsPersistentCacheDir(cm.Target, cm.LLBCacheMountSharingMode()),
		),
	}, nil
}

func (cm *CacheMount) LLBCacheMountSharingMode() llb.CacheMountSharingMode {
	switch cm.Access {
	case CachePrivate:
		return llb.CacheMountPrivate
	case CacheLocked:
		return llb.CacheMountLocked
	}

	return llb.CacheMountShared
}
