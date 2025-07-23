package configdomain

type LineageDisplayLocation string

const (
	LineageDisplayLocationNone            LineageDisplayLocation = "none"
	LineageDisplayLocationProposalComment LineageDisplayLocation = "comment"
	LineageDisplayLocationProposalBody    LineageDisplayLocation = "body"
	LineageDisplayLocationTerminal        LineageDisplayLocation = "terminal"
)

func (self LineageDisplayLocation) String() string {
	return string(self)
}
