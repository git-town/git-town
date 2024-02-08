package main

import (
	"sort"

	"golang.org/x/exp/maps"
)

type UserCollector struct {
	users map[string]struct{}
}

func NewUserCollector() UserCollector {
	return UserCollector{
		users: map[string]struct{}{},
	}
}

func (self *UserCollector) AddUser(id string) {
	self.users[id] = struct{}{}
}

func (self *UserCollector) AddUsers(users UserCollector) {
	for _, user := range users.Users() {
		self.AddUser(user)
	}
}

func (self *UserCollector) Users() []string {
	result := maps.Keys(self.users)
	sort.Strings(result)
	return result
}
