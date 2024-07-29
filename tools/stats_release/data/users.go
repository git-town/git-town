package data

import (
	"github.com/git-town/git-town/v14/src/gohacks"
)

// collection of unique GitHub usernames
type Users struct {
	gohacks.Set[string]
}

func NewUsers(users ...string) Users {
	result := Users{gohacks.NewSet[string](users...)}
	return result
}

func (self Users) AddUsers(other Users) {
	self.Add(other.Values()...)
}
