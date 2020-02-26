package selective

type T__ struct{}

func First(t ...T__) T__ {
	return t[0]
}

func Last(t ...T__) T__ {
	return t[len(t)-1]
}
