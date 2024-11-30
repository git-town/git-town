package configdomain

import "github.com/git-town/git-town/v16/internal/git/gitdomain"

type LineageBuilder struct {
	data map[gitdomain.LocalBranchName]gitdomain.LocalBranchName
}

func NewLineageBuilder() LineageBuilder {
	data := new(map[gitdomain.LocalBranchName]gitdomain.LocalBranchName)
	return LineageBuilder{
		data: *data,
	}
}

func (self LineageBuilder) Add(branch, parent gitdomain.LocalBranchName) LineageBuilder {
	self.data[branch] = parent
	return self
}

func (self *LineageBuilder) Lineage() Lineage {
	return Lineage{
		data: self.data,
	}
}
