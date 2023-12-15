package configdomain

// GiteaToken is a bearer token to use with the Gitea API.
type GiteaToken string

func (self GiteaToken) String() string {
	return string(self)
}
