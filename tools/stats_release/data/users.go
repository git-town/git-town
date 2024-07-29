package data

import (
	"github.com/git-town/git-town/v14/src/gohacks"
)

// collection of unique GitHub usernames
type Users struct {
	gohacks.Set[string]
}

func NewUsers(users ...string) Users {
	return Users{Set: gohacks.NewSet(users...)}
}

func (self Users) AddUsers(other Users) {
	self.AddMany(other.Values()...)
}
