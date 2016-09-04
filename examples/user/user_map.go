package iterator

type UserSlice []User

func (t UserSlice) Map(fn func(User) User) UserSlice {
	o := make(UserSlice, len(t))
	for i := range t {
		o[i] = fn(t[i])
	}
	return o
}
