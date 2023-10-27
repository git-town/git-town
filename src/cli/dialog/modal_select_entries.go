package dialog

// ModalSelectEntries is a collection of ModalEntry.
type ModalSelectEntries []ModalSelectEntry

// IndexOfValue provides the index of the entry with the given value,
// or nil if the given value is not in the list.
func (self ModalSelectEntries) IndexOfValue(value string) *int {
	for e, entry := range self {
		if entry.Value == value {
			return &e
		}
	}
	return nil
}
