package opcodes

// ProgramEndOfBranch marks the end of the program to sync a branch.
// All opcodes after this opcode are not for syncing this branch.
// They might sync another branch, or do something else.
type ProgramEndOfBranch struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}
