package dialog

// modalSelectStatus represents the different states that a modalSelect instance can be in.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type modalSelectStatus struct {
	name string
}

var (
	modalSelectStatusNew       = modalSelectStatus{"new"}       //nolint:gochecknoglobals
	modalSelectStatusSelecting = modalSelectStatus{"selecting"} //nolint:gochecknoglobals
	modalSelectStatusSelected  = modalSelectStatus{"selected"}  //nolint:gochecknoglobals
	modalSelectStatusAborted   = modalSelectStatus{"aborted"}   //nolint:gochecknoglobals
)
