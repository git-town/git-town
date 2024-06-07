package gohacks

type MutableList[E any, L ~[]E] struct {
	list *L
}

func (self MutableList[E, L]) Append(element E) {
	*self.list = append(*self.list, element)
}

func (self MutableList[E, L]) List() L {
	if self.list == nil {
		return make(L, 0)
	}
	return *self.list
}
