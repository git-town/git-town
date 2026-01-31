package bytestream

import "bytes"

// NullDelineated is raw Git output that is delineated by null bytes.
type NullDelineated []byte

func (self NullDelineated) ToNewlines() NewlineDelineated {
	return bytes.ReplaceAll(self, []byte{0x00}, []byte{'\n', '\n'})
}
