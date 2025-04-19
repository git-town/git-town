package gitlab

import (
	"net/http"
	"time"
)

type (
	BulkImportsServiceInterface interface {
		StartMigration(startMigrationOptions *BulkImportStartMigrationOptions, options ...RequestOptionFunc) (*BulkImportStartMigrationResponse, *Response, error)
	}

	// BulkImportsService handles communication with GitLab's direct transfer API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/bulk_imports/
	BulkImportsService struct {
		client *Client
	}
)

var _ BulkImportsServiceInterface = (*BulkImportsService)(nil)

// BulkImportStartMigrationConfiguration represents the available configuration options to start a migration.
//
// GitLab API docs: https://docs.gitlab.com/api/bulk_imports/#start-a-new-group-or-project-migration
type BulkImportStartMigrationConfiguration struct {
	URL         *string `json:"url,omitempty"`
	AccessToken *string `json:"access_token,omitempty"`
}

// BulkImportStartMigrationEntity represents the available entity options to start a migration.
//
// GitLab API docs: https://docs.gitlab.com/api/bulk_imports/#start-a-new-group-or-project-migration
type BulkImportStartMigrationEntity struct {
	SourceType           *string `json:"source_type,omitempty"`
	SourceFullPath       *string `json:"source_full_path,omitempty"`
	DestinationSlug      *string `json:"destination_slug,omitempty"`
	DestinationNamespace *string `json:"destination_namespace,omitempty"`
	MigrateProjects      *bool   `json:"migrate_projects,omitempty"`
	MigrateMemberships   *bool   `json:"migrate_memberships,omitempty"`
}

// BulkImportStartMigrationOptions represents the available start migration options.
//
// GitLab API docs: https://docs.gitlab.com/api/bulk_imports/#start-a-new-group-or-project-migration
type BulkImportStartMigrationOptions struct {
	Configuration *BulkImportStartMigrationConfiguration `json:"configuration,omitempty"`
	Entities      []BulkImportStartMigrationEntity       `json:"entities,omitempty"`
}

// BulkImportStartMigrationResponse represents the start migration response.
//
// GitLab API docs: https://docs.gitlab.com/api/bulk_imports/#start-a-new-group-or-project-migration
type BulkImportStartMigrationResponse struct {
	ID          int       `json:"id"`
	Status      string    `json:"status"`
	SourceType  string    `json:"source_type"`
	SourceURL   string    `json:"source_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	HasFailures bool      `json:"has_failures"`
}

// StartMigration starts a migration.
//
// GitLab API docs: https://docs.gitlab.com/api/bulk_imports/#start-a-new-group-or-project-migration
func (b *BulkImportsService) StartMigration(startMigrationOptions *BulkImportStartMigrationOptions, options ...RequestOptionFunc) (*BulkImportStartMigrationResponse, *Response, error) {
	request, err := b.client.NewRequest(http.MethodPost, "bulk_imports", startMigrationOptions, options)
	if err != nil {
		return nil, nil, err
	}

	startMigrationResponse := new(BulkImportStartMigrationResponse)
	response, err := b.client.Do(request, startMigrationResponse)
	if err != nil {
		return nil, response, err
	}

	return startMigrationResponse, response, nil
}
