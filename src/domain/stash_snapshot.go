package domain

// StashSnapshot is a snapshot of th state of Git stash at a given point in time.
type StashSnapshot struct {
	Amount int // the amount of Git stash entries
}

func EmptyStashSnapshot() StashSnapshot {
	return StashSnapshot{Amount: 0}
}
