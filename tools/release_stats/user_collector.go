package main

import (
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

// a list of unique GitHub usernames
type Users struct {
	users map[string]struct{}
}

func NewUsers() Users {
	return Users{
		users: map[string]struct{}{},
	}
}

func (self *Users) AddUser(id string) {
	self.users[id] = struct{}{}
}

func (self *Users) AddUsers(users Users) {
	for _, user := range users.Users() {
		self.AddUser(user)
	}
}

func (self *Users) Users() []string {
	result := maps.Keys(self.users)
	sort.Slice(result, func(i, j int) bool { return strings.ToLower(result[i]) < strings.ToLower(result[j]) })
	return result
}
