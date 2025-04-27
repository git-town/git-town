package forgedomain

type ProposalCommentSearchOptions struct {
}

type ProposalCommentSortBy string

const (
	ProposalCommentSortByCreatedAt ProposalCommentSortBy = "created_at"
	ProposalCommentSortByUpdatedAt ProposalCommentSortBy = "updated_at"
)

func (self ProposalCommentSortBy) String() string {
	return string(self)
}

type ProposalCommentQueryOptions struct {
	limit  int
	sortBy ProposalCommentSortBy
}

func (self *ProposalCommentQueryOptions) Limit() int {
	return self.limit
}

func (self *ProposalCommentQueryOptions) SortBy() ProposalCommentSortBy {
	return self.sortBy
}

func NewProposalCommentQueryOptions() *ProposalCommentQueryOptions {
	return &ProposalCommentQueryOptions{
		limit:  100,
		sortBy: ProposalCommentSortByCreatedAt,
	}
}

type ConfigureProposalCommentQueryOptions func(options *ProposalCommentQueryOptions)

func WithLimit(limit int) ConfigureProposalCommentQueryOptions {
	return func(options *ProposalCommentQueryOptions) {
		options.limit = limit
	}
}

func WithSortBy(sortBy ProposalCommentSortBy) ConfigureProposalCommentQueryOptions {
	return func(options *ProposalCommentQueryOptions) {
		options.sortBy = sortBy
	}
}
