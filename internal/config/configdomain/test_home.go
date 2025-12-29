package configdomain

// TestHome is the path of the home directory when running inside a test environment.
type TestHome string

func (self TestHome) String() string {
	return string(self)
}
