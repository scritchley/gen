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
