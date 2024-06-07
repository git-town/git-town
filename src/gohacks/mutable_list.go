package gohacks

type MutableList[E any, L ~[]E] struct {
	list *L
}

func (self *MutableList[E, L]) Append(element E) {
	self.initializeIfNeeded()
	*self.list = append(*self.list, element)
}

func (self *MutableList[E, L]) List() L {
	self.initializeIfNeeded()
	return *self.list
}

// indicates whether this list has the non-functional zero value
func (self *MutableList[E, L]) isInitialized() bool {
	return self.list != nil
}

// initializes this list to an empty functioning state
func (self *MutableList[E, L]) initialize() {
	list := make(L, 0)
	self.list = &list
}

// get this list to an empty functioning state if it isn't in one yet
func (self *MutableList[E, L]) initializeIfNeeded() {
	if !self.isInitialized() {
		self.initialize()
	}
}
