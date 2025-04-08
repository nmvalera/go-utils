package common

import "reflect"

// Ptr returns the pointer to the value passed in.
func Ptr[T any](v T) *T {
	return &v
}

// PtrSlice returns a pointer to a slice of pointers to all of the provided values.
func PtrSlice[T any](v ...T) *[]*T {
	slice := make([]*T, len(v))
	for i := range v {
		slice[i] = Ptr(v[i])
	}
	return &slice
}

// Val returns the value pointed to by the pointer passed in.
// If the pointer is nil, it returns the zero value of the type.
func Val[T any](v *T) T {
	if v == nil {
		return reflect.Zero(reflect.TypeOf(v).Elem()).Interface().(T)
	}
	return *v
}

func ValSlice[T any](v ...*T) []T {
	values := make([]T, len(v))
	for i := range v {
		values[i] = Val(v[i])
	}
	return values
}

func Copy[T any](v *T) *T {
	if v == nil {
		return nil
	}
	return Ptr(Val(v))
}

func CopySlice[T any](v ...*T) []*T {
	slice := make([]*T, len(v))
	for i := range v {
		slice[i] = Copy(v[i])
	}
	return slice
}
