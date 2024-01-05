package dialog

// modalSelectStatus represents the different states that a modalSelect instance can be in.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type modalSelectStatus string

const (
	modalSelectStatusNew       = modalSelectStatus("new")
	modalSelectStatusSelecting = modalSelectStatus("selecting")
	modalSelectStatusSelected  = modalSelectStatus("selected")
	modalSelectStatusAborted   = modalSelectStatus("aborted")
)
