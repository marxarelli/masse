package state

type Git struct {
	Repo    string       `json:"git"`
	Ref     string       `json:"ref"`
	Options []*GitOption `json:"options"`
}

type GitOption struct {
	*KeepGitDir
}

type KeepGitDir struct {
	KeepGitDir bool `json:"keepGitDir"`
}
