package iterator

type Type interface{}

// TypeSlice is a slice of Type.
type TypeSlice []Type

// Map calls the provided func for each element in t and returns a new TypeSlice.
func (t TypeSlice) Map(fn func(Type) Type) TypeSlice {
	o := make(TypeSlice, len(t))
	for i := range t {
		o[i] = fn(t[i])
	}
	return o
}

// GroupBy returns a map of string keys to TypeSlice using the provided groupBy fn.
func (t TypeSlice) GroupBy(groupBy func(Type) string) map[string]TypeSlice {
	groups := make(map[string]TypeSlice)
	for i := range t {
		group := groupBy(t[i])
		groups[group] = append(groups[group], t[i])
	}
	return groups
}

type TypeAccumulator interface{}

func (t TypeSlice) Reduce(accumulator func(TypeAccumulator, Type) TypeAccumulator, initial TypeAccumulator) TypeAccumulator {
	for i := range t {
		initial = accumulator(initial, t[i])
	}
	return initial
}

func (t TypeSlice) ignoreThisMethod() {

}
