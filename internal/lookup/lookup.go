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
