package domain

// StashSnapshot is a snapshot of the state of Git stash at a given point in time.
type StashSnapshot int

func EmptyStashSnapshot() StashSnapshot {
	return 0
}
