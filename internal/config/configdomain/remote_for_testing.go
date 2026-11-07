package configdomain

type RemoteForTesting string

func (self RemoteForTesting) String() string {
	return string(self)
}
