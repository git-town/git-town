package commitmessage

import "strings"

// Split separates the parts of the given commit message.
func Split(message string) Parts {
	parts := strings.SplitN(message, "\n", 2)
	title := parts[0]
	body := ""
	if len(parts) == 2 {
		body = parts[1]
	}
	for strings.HasPrefix(body, "\n") {
		body = body[1:]
	}
	return Parts{
		Title: title,
		Body:  body,
	}
}
