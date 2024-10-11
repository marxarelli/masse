package v1

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/pkg/errors"
)

func vError(v cue.Value, err error) error {
	if err != nil {
		return errorf(v, "compile error: %s", err)
	}
	return nil
}

func errorf(v cue.Value, msg string, args ...any) error {
	return errors.Errorf(fmt.Sprintf("%s: %s at %s", v.Path(), msg, v.Pos()), args...)
}
