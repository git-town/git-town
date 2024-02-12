package data

import (
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

// counts users and their contributions
type ContributionCounter struct {
	list map[string]int
}

func NewContributionCounter() ContributionCounter {
	return ContributionCounter{
		list: map[string]int{},
	}
}

func (self *ContributionCounter) AddUser(id string) {
	self.list[id] = self.list[id] + 1
}

func (self *ContributionCounter) AddUsers(users ContributionCounter) {
	for _, contributor := range users.Contributors() {
		for i := 0; i < contributor.ContributionCount; i++ {
			self.AddUser(contributor.Username)
		}
	}
}

func (self *ContributionCounter) Contributors() []Contributor {
	usernames := maps.Keys(self.list)
	sort.Slice(usernames, func(i, j int) bool { return strings.ToLower(usernames[i]) < strings.ToLower(usernames[j]) })
	result := make([]Contributor, len(self.list))
	for u, username := range usernames {
		result[u] = Contributor{
			ContributionCount: self.list[username],
			Username:          username,
		}
	}
	return result
}
