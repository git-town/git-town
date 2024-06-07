package gohacks

type MutableList[E any, L ~[]E] struct {
	list *L
}

func NewMutableList[E any, L ~[]E]() MutableList[E, L] {
	list := make(L, 0)
	return MutableList[E, L]{list: &list}
}

func (self MutableList[E, L]) Append(element E) {
	*self.list = append(*self.list, element)
}

func (self MutableList[E, L]) List() L {
	return *self.list
}
