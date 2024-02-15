package opcodes

// SkipCurrentBranch is a mock opcode to be used instead of
// running another program.
// This is used when ignoring the remaining opcodes for a particular branch.
type SkipCurrentBranch struct {
	undeclaredOpcodeMethods
}
