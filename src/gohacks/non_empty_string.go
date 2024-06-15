package gohacks

type NonEmptyString string

func NewNonEmptyString(value string) NonEmptyString {
	if len(value) == 0 {
		panic("provided an empty value for a non-empty string")
	}
	return NonEmptyString(value)
}
