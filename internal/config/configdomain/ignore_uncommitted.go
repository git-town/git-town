package configdomain

import "strconv"

// IgnoreUncommitted indicates whether to allow uncommitted changes when shipping.
type IgnoreUncommitted bool

func (self IgnoreUncommitted) AllowUncommitted() bool {
	return bool(self)
}

func (self IgnoreUncommitted) DisAllowUncommitted() bool {
	return !self.AllowUncommitted()
}

func (self IgnoreUncommitted) String() string {
	return strconv.FormatBool(bool(self))
}
