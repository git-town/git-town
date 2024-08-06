package data

import (
	"github.com/git-town/git-town/v14/pkg/set"
)

// collection of unique GitHub usernames
type Users struct {
	set.Set[string]
}

func NewUsers(users ...string) Users {
	return Users{set.NewSet(users...)}
}

func (self Users) AddUsers(other Users) {
	self.Add(other.Values()...)
}
