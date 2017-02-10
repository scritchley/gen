package user

// User is a user model.
//go:generate gen -src github.com/scritchley/gen/examples/iterator -dest User
type User struct {
	Name string
	Age  int
}
