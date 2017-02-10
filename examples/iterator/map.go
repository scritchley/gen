package iterator

// T__Slice is a slice of T__.
type T__Slice []T__

// Map calls the provided func for each element in t and returns a new T__Slice.
func (t T__Slice) Map(fn func(T__) T__) T__Slice {
	o := make(T__Slice, len(t))
	for i := range t {
		o[i] = fn(t[i])
	}
	return o
}
