package v1

import (
	"time"

	"github.com/containerd/containerd/platforms"
	protoyaml "gitlab.wikimedia.org/dduvall/protoyaml/protoyaml"
)

func (platform *Platform) UnmarshalProtoYAML(node *protoyaml.Node) error {
	return node.WithEntry("name", func(node *protoyaml.Node) error {
		var name string
		err := node.Decode(&name)
		if err != nil {
			return err
		}

		p, err := platforms.Parse(name)
		if err != nil {
			return err
		}

		p = platforms.Normalize(p)

		platform.Os = p.OS
		platform.Architecture = p.Architecture
		platform.Variant = p.Variant

		return nil
	})
}

func (env *Env) UnmarshalProtoYAML(node *protoyaml.Node) error {
	return node.WithEntry("variables", func(node *protoyaml.Node) error {
		env.Variables = make([]*NameValue, node.Len())

		return node.WithPairs(func(i int, key *protoyaml.Node, value *protoyaml.Node) error {
			nv := &NameValue{}

			err := key.Decode(&nv.Name)
			if err != nil {
				return err
			}

			err = value.Decode(&nv.Value)
			if err != nil {
				return err
			}

			env.Variables[i] = nv

			return nil
		})
	})
}

func (ts *Timestamp) UnmarshalProtoYAML(node *protoyaml.Node) error {
	return node.WithEntry("seconds", func(node *protoyaml.Node) error {
		t, err := time.Parse(time.RFC3339, node.Value)

		if err != nil {
			return node.Errorf("failed to parse rfc3339 time value %v", node.Value)
		}

		ts.Seconds = t.Unix()

		return nil
	})

	return nil
}
