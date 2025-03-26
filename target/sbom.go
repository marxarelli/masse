package target

import (
	"cuelang.org/go/cue"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type SBOM struct {
	Generator  string     `json:"generator"`
	Parameters common.Env `json:"parameters"`
	Scan       cue.Value  `json:"scan"`
}

func (sbom SBOM) ScanAll() bool {
	scan, _ := sbom.Scan.String()
	return scan == "all"
}

func (sbom SBOM) ScanChainRefs() []string {
	refs := []string{}

	iter, err := sbom.Scan.List()
	if err != nil {
		return refs
	}

	for iter.Next() {
		ref, err := iter.Value().String()
		if err == nil {
			refs = append(refs, ref)
		}
	}

	return refs
}

func (sbom SBOM) ScanChainRefMap() map[string]bool {
	refs := sbom.ScanChainRefs()
	refMap := make(map[string]bool, len(refs))

	for _, ref := range refs {
		refMap[ref] = true
	}

	return refMap
}
