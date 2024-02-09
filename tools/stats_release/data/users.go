package data

import (
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

// collection of unique GitHub usernames
type Users struct {
	list map[string]struct{}
}

func NewUsers() Users {
	return Users{
		list: map[string]struct{}{},
	}
}

func (self *Users) AddUser(id string) {
	self.list[id] = struct{}{}
}

func (self *Users) AddUsers(users Users) {
	for _, user := range users.Users() {
		self.AddUser(user)
	}
}

func (self *Users) Users() []string {
	result := maps.Keys(self.list)
	sort.Slice(result, func(i, j int) bool { return strings.ToLower(result[i]) < strings.ToLower(result[j]) })
	return result
}
