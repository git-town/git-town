package messages

import (
	"strconv"

	"github.com/google/uuid"
)

type IdGenerator interface {
	newId() func() string
}

type Incrementing struct {
	next int
}

func (self *Incrementing) NewId() string {
	result := strconv.Itoa(self.next)
	self.next++
	return result
}

type UUID struct {
	next int
}

func (i UUID) NewId() string {
	return uuid.New().String()
}
