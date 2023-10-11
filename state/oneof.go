package state

import "reflect"

func oneof[T any](xs any) (T, bool) {
	var zero T

	xv := reflect.ValueOf(xs)
	if xv.Kind() == reflect.Pointer {
		xv = xv.Elem()
	}

	xt := xv.Type()
	nf := xv.NumField()

	for i := 0; i < nf; i++ {
		val := xv.Field(i)

		if xt.Field(i).Anonymous {
			isNil := false

			switch val.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
				isNil = val.IsNil()
			}

			if !isNil {
				v, ok := val.Interface().(T)
				if ok {
					return v, true
				}
			}
		}
	}

	return zero, false
}
