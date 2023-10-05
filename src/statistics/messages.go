package statistics

import "fmt"

// Messages is a statistics implementation that collects only messages.
type Messages struct {
	messages []string
}

func (m *Messages) Queue(message string) {
	m.messages = append(m.messages, message)
}

func (m *Messages) PrintQueue() {
	for _, message := range m.messages {
		fmt.Println(message)
	}
}
