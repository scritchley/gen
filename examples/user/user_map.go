package user

// UserSlice is a slice of User.
type UserSlice []User

// Map calls the provided func for each element in t and returns a new UserSlice.
func (t UserSlice) Map(fn func(User) User) UserSlice {
	o := make(UserSlice, len(t))
	for i := range t {
		o[i] = fn(t[i])
	}
	return o
}

// GroupBy returns a map of string keys to UserSlice using the provided groupBy fn.
func (t UserSlice) GroupBy(groupBy func(User) string) map[string]UserSlice {
	groups := make(map[string]UserSlice)
	for i := range t {
		group := groupBy(t[i])
		groups[group] = append(groups[group], t[i])
	}
	return groups
}

type UserAccumulator interface{}

func (t UserSlice) Reduce(accumulator func(UserAccumulator, User) UserAccumulator, initial UserAccumulator) UserAccumulator {
	for i := range t {
		initial = accumulator(initial, t[i])
	}
	return initial
}
