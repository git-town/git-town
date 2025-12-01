package configdomain

type Browser string

func (self Browser) String() string {
	return string(self)
}

func (self Browser) NoBrowser() bool {
	return self == NoBrowser
}

const NoBrowser = "(none)"
