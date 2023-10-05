package state

import "gitlab.wikimedia.org/dduvall/phyton/common"

type Copy struct {
	Source      []common.Glob `json:"source"`
	From        ChainRef      `json:"from"`
	Destination string        `json:"destination"`
	Options     []*CopyOption `json:"options"`
}

type CopyOption struct {
	*common.Creation      `json:",inline"`
	*common.User          `json:",inline"`
	*common.Group         `json:",inline"`
	*common.Mode          `json:",inline"`
	*common.Include       `json:",inline"`
	*common.Exclude       `json:",inline"`
	*FollowSymlinks       `json:",inline"`
	*CopyDirectoryContent `json:",inline"`
}

type FollowSymlinks struct {
	FollowSymlinks bool `json:"followSymlinks"`
}

type CopyDirectoryContent struct {
	CopyDirectoryContent bool `json:"copyDirectoryContent"`
}
