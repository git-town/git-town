package statistics

import (
	"encoding/json"
	"fmt"
)

// None is a statistics implementation that does nothing.
type Messages struct {
	list []string
}

// NewMessages provides new instances of Messages for testing only!
func NewMessages(messages ...string) Messages {
	return Messages{list: messages}
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (m Messages) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.list)
}

func (m *Messages) Queue(message string) {
	m.list = append(m.list, message)
}

func (m *Messages) PrintAll() {
	for _, message := range m.list {
		fmt.Println("\n" + message)
	}
}

// UnmarshalJSON is used when de-serializing JSON into a Location.
func (m *Messages) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.list)
}
