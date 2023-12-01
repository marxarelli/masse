package common

import (
	"fmt"
	"sort"
)

type Env map[string]string

func (env *Env) Sort() []string {
	names := make([]string, len(*env))

	i := 0
	for name := range *env {
		names[i] = name
		i++
	}

	sort.Strings(names)

	return names
}

func (env *Env) Assignments() []string {
	assigns := make([]string, len(*env))

	for i, name := range env.Sort() {
		assigns[i] = fmt.Sprintf("%s=%s", name, (*env)[name])
	}

	return assigns
}
