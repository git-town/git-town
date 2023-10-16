package common

import "strings"

func ParseCommitMessage(message string) (title, body string) {
	parts := strings.SplitN(message, "\n", 2)
	title = parts[0]
	if len(parts) == 2 {
		body = parts[1]
	} else {
		body = ""
	}
	for strings.HasPrefix(body, "\n") {
		body = body[1:]
	}
	return
}
