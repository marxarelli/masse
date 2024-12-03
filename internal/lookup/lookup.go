package lookup

import (
	"cuelang.org/go/cue"
	"github.com/pkg/errors"
)

func Lookup(root cue.Value, path string) cue.Value {
	return root.LookupPath(cue.ParsePath(path))
}

func Existing(root cue.Value, path string) (cue.Value, error) {
	v := root.LookupPath(cue.ParsePath(path))
	if err := v.Err(); err != nil {
		return v, err
	}

	if !v.Exists() {
		return v, errors.Errorf("%s not found", path)
	}

	return v, nil
}

func String(root cue.Value, path string) (string, error) {
	v, err := Existing(root, path)
	if err != nil {
		return "", err
	}

	return v.String()
}

func Bool(root cue.Value, path string) (bool, error) {
	v, err := Existing(root, path)
	if err != nil {
		return true, err
	}

	return v.Bool()
}

func Int64(root cue.Value, path string) (int64, error) {
	v, err := Existing(root, path)
	if err != nil {
		return int64(0), err
	}

	return v.Int64()
}

func Bytes(root cue.Value, path string) ([]byte, error) {
	v, err := Existing(root, path)
	if err != nil {
		return nil, err
	}

	return v.Bytes()
}

func Each(root cue.Value, path string, f func(cue.Value) error) error {
	v, err := Existing(root, path)
	if err != nil {
		return err
	}

	iter, err := v.List()
	if err != nil {
		return err
	}

	for iter.Next() {
		entry := iter.Value()
		if err := entry.Err(); err != nil {
			return err
		}

		err := f(entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func EachField(root cue.Value, path string, f func(label string, value cue.Value) error) error {
	v, err := Existing(root, path)
	if err != nil {
		return err
	}

	iter, err := v.Fields()
	if err != nil {
		return err
	}

	for iter.Next() {
		sel := iter.Selector()

		if sel.LabelType() == cue.StringLabel {
			label := sel.String()
			value := iter.Value()

			if err := value.Err(); err != nil {
				return err
			}

			err := f(label, value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func EachOrValue(root cue.Value, path string, f func(cue.Value) error) error {
	v, err := Existing(root, path)
	if err != nil {
		return err
	}

	if v.Kind() != cue.ListKind {
		return f(v)
	}

	return Each(root, path, f)
}

func DecodeListOrSingle[S ~[]E, E comparable](root cue.Value, path string) (S, error) {
	var nilE E

	s := S{}
	v := Lookup(root, path)

	if !v.Exists() {
		return s, nil
	}

	err := v.Err()
	if err != nil {
		return s, err
	}

	if !v.IsConcrete() {
		return s, errors.Errorf("cannot decode non-concrete value %q", v)
	}

	if v.Kind() != cue.ListKind {
		var e E
		err := v.Decode(&e)
		if err != nil {
			return s, err
		}

		if e == nilE {
			return S{}, nil
		}

		return S{e}, nil
	}

	iter, err := v.List()
	if err != nil {
		return s, err
	}

	for iter.Next() {
		var e E
		item := iter.Value()
		err := item.Decode(&e)
		if err != nil {
			return s, err
		}

		if e != nilE {
			s = append(s, e)
		}
	}

	return s, nil
}

func DecodeOptions[S ~[]E, E comparable](root cue.Value) (S, error) {
	defaults := S{}
	var err error

	if Lookup(root, "#defaultOptions").Exists() {
		defaults, err = DecodeListOrSingle[S](root, "#defaultOptions")
		if err != nil {
			return defaults, err
		}
	}

	options, err := DecodeListOrSingle[S](root, "options")
	if err != nil {
		return defaults, err
	}

	return append(defaults, options...), nil
}

func WithDiscriminatorField[T ~string, R any](v cue.Value, f func(T) (R, bool, error)) (R, bool, error) {
	var nilR R
	iter, err := v.Fields(cue.Final(), cue.Concrete(true), cue.Optional(false))
	if err != nil {
		return nilR, false, err
	}

	for iter.Next() {
		sel := iter.Selector()

		if sel.LabelType() == cue.StringLabel {
			r, matched, err := f(T(sel.String()))
			if matched {
				return r, true, err
			}
		}
	}

	return nilR, false, nil
}
