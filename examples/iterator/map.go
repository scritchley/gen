package iterator

type T__Slice []T__

func (t T__Slice) Map(fn func(T__) T__) T__Slice {
	o := make(T__Slice, len(t))
	for i := range t {
		o[i] = fn(t[i])
	}
	return o
}
