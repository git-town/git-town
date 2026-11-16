package pageant

import (
	"io"

	"golang.org/x/crypto/ssh/agent"
)

// New returns a new ssh-agent client (see golang.org/x/crypto/ssh/agent)
// connected to Pageant.
func New() agent.Agent {
	return agent.NewClient(&conn{})
}

// conn has no lock of its own: the agent client serializes all calls.
type conn struct {
	buf []byte
}

func (c *conn) Write(p []byte) (int, error) {
	resp, err := query(p)
	if err != nil {
		return 0, err
	}
	// Replace, not append: leftovers of a partly read previous response
	// must not leak into this one.
	c.buf = resp
	return len(p), nil
}

func (c *conn) Read(p []byte) (int, error) {
	if len(c.buf) == 0 {
		return 0, io.EOF
	}
	n := copy(p, c.buf)
	c.buf = c.buf[n:]
	return n, nil
}
