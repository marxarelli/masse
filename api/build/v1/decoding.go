package v1

import (
	"time"

	protoyaml "gitlab.wikimedia.org/dduvall/protoyaml/protoyaml"
)

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
