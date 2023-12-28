package configdomain

import (
	"golang.org/x/exp/maps"
)

// SingleCache caches a single Git configuration type (local or global).
type SingleCache map[Key]string

// Clone provides a copy of this GitConfiguration instance.
func (self SingleCache) Clone() SingleCache {
	result := SingleCache{}
	maps.Copy(result, self)
	return result
}
