package commitmessage

import "strings"

type Parts struct {
	Title string
	Body  string
}

// Split splits the given commit message into its header and body parts.
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
