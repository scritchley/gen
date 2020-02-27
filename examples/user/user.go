package user

// User is a user model.
//go:generate gen -src github.com/scritchley/gen/examples/iterator -replace Type=User -exclude Type,TypeSlice.ignoreThisMethod
type User struct {
	Name string
	Age  int
}

type UserByAgeCount map[string]int
