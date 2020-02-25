package iterator

type T__ interface{}

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

// GroupBy returns a map of string keys to T__Slice using the provided groupBy fn.
func (t T__Slice) GroupBy(groupBy func(T__) string) map[string]T__Slice {
	groups := make(map[string]T__Slice)
	for i := range t {
		group := groupBy(t[i])
		groups[group] = append(groups[group], t[i])
	}
	return groups
}

type T__Accumulator interface{}

func (t T__Slice) Reduce(accumulator func(T__Accumulator, T__) T__Accumulator, initial T__Accumulator) T__Accumulator {
	for i := range t {
		initial = accumulator(initial, t[i])
	}
	return initial
}
