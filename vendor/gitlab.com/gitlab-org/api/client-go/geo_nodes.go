//
// Copyright 2021, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import "net/http"

type (
	// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
	GeoNodesServiceInterface interface {
		// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
		CreateGeoNode(*CreateGeoNodesOptions, ...RequestOptionFunc) (*GeoNode, *Response, error)
		// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
		ListGeoNodes(*ListGeoNodesOptions, ...RequestOptionFunc) ([]*GeoNode, *Response, error)
		// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
		GetGeoNode(int64, ...RequestOptionFunc) (*GeoNode, *Response, error)
		// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
		EditGeoNode(int64, *UpdateGeoNodesOptions, ...RequestOptionFunc) (*GeoNode, *Response, error)
		// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
		DeleteGeoNode(int64, ...RequestOptionFunc) (*Response, error)
		// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
		RepairGeoNode(int64, ...RequestOptionFunc) (*GeoNode, *Response, error)
		// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
		RetrieveStatusOfAllGeoNodes(...RequestOptionFunc) ([]*GeoNodeStatus, *Response, error)
		// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
		RetrieveStatusOfGeoNode(int64, ...RequestOptionFunc) (*GeoNodeStatus, *Response, error)
	}

	// GeoNodesService handles communication with Geo Nodes related methods
	// of GitLab API.
	// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
	//
	// GitLab API docs: https://docs.gitlab.com/api/geo_nodes/
	GeoNodesService struct {
		client *Client
	}
)

// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
var _ GeoNodesServiceInterface = (*GeoNodesService)(nil)

// GeoNode represents a GitLab Geo Node.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs: https://docs.gitlab.com/api/geo_nodes/
type GeoNode struct {
	ID                               int64        `json:"id"`
	Name                             string       `json:"name"`
	URL                              string       `json:"url"`
	InternalURL                      string       `json:"internal_url"`
	Primary                          bool         `json:"primary"`
	Enabled                          bool         `json:"enabled"`
	Current                          bool         `json:"current"`
	FilesMaxCapacity                 int64        `json:"files_max_capacity"`
	ReposMaxCapacity                 int64        `json:"repos_max_capacity"`
	VerificationMaxCapacity          int64        `json:"verification_max_capacity"`
	SelectiveSyncType                string       `json:"selective_sync_type"`
	SelectiveSyncShards              []string     `json:"selective_sync_shards"`
	SelectiveSyncNamespaceIDs        []int64      `json:"selective_sync_namespace_ids"`
	MinimumReverificationInterval    int64        `json:"minimum_reverification_interval"`
	ContainerRepositoriesMaxCapacity int64        `json:"container_repositories_max_capacity"`
	SyncObjectStorage                bool         `json:"sync_object_storage"`
	CloneProtocol                    string       `json:"clone_protocol"`
	WebEditURL                       string       `json:"web_edit_url"`
	WebGeoProjectsURL                string       `json:"web_geo_projects_url"`
	Links                            GeoNodeLinks `json:"_links"`
}

// GeoNodeLinks represents links for GitLab GeoNode.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs: https://docs.gitlab.com/api/geo_nodes/
type GeoNodeLinks struct {
	Self   string `json:"self"`
	Status string `json:"status"`
	Repair string `json:"repair"`
}

// CreateGeoNodesOptions represents the available CreateGeoNode() options.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#create-a-new-geo-node
type CreateGeoNodesOptions struct {
	Primary                          *bool     `url:"primary,omitempty" json:"primary,omitempty"`
	Enabled                          *bool     `url:"enabled,omitempty" json:"enabled,omitempty"`
	Name                             *string   `url:"name,omitempty" json:"name,omitempty"`
	URL                              *string   `url:"url,omitempty" json:"url,omitempty"`
	InternalURL                      *string   `url:"internal_url,omitempty" json:"internal_url,omitempty"`
	FilesMaxCapacity                 *int64    `url:"files_max_capacity,omitempty" json:"files_max_capacity,omitempty"`
	ReposMaxCapacity                 *int64    `url:"repos_max_capacity,omitempty" json:"repos_max_capacity,omitempty"`
	VerificationMaxCapacity          *int64    `url:"verification_max_capacity,omitempty" json:"verification_max_capacity,omitempty"`
	ContainerRepositoriesMaxCapacity *int64    `url:"container_repositories_max_capacity,omitempty" json:"container_repositories_max_capacity,omitempty"`
	SyncObjectStorage                *bool     `url:"sync_object_storage,omitempty" json:"sync_object_storage,omitempty"`
	SelectiveSyncType                *string   `url:"selective_sync_type,omitempty" json:"selective_sync_type,omitempty"`
	SelectiveSyncShards              *[]string `url:"selective_sync_shards,omitempty" json:"selective_sync_shards,omitempty"`
	SelectiveSyncNamespaceIDs        *[]int64  `url:"selective_sync_namespace_ids,omitempty" json:"selective_sync_namespace_ids,omitempty"`
	MinimumReverificationInterval    *int64    `url:"minimum_reverification_interval,omitempty" json:"minimum_reverification_interval,omitempty"`
}

// CreateGeoNode creates a new Geo Node.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#create-a-new-geo-node
func (s *GeoNodesService) CreateGeoNode(opt *CreateGeoNodesOptions, options ...RequestOptionFunc) (*GeoNode, *Response, error) {
	return do[*GeoNode](s.client,
		withMethod(http.MethodPost),
		withPath("geo_nodes"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListGeoNodesOptions represents the available ListGeoNodes() options.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#retrieve-configuration-about-all-geo-nodes
type ListGeoNodesOptions struct {
	ListOptions
}

// ListGeoNodes gets a list of geo nodes.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#retrieve-configuration-about-all-geo-nodes
func (s *GeoNodesService) ListGeoNodes(opt *ListGeoNodesOptions, options ...RequestOptionFunc) ([]*GeoNode, *Response, error) {
	return do[[]*GeoNode](s.client,
		withPath("geo_nodes"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetGeoNode gets a specific geo node.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#retrieve-configuration-about-a-specific-geo-node
func (s *GeoNodesService) GetGeoNode(id int64, options ...RequestOptionFunc) (*GeoNode, *Response, error) {
	return do[*GeoNode](s.client,
		withPath("geo_nodes/%d", id),
		withRequestOpts(options...),
	)
}

// UpdateGeoNodesOptions represents the available EditGeoNode() options.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#edit-a-geo-node
type UpdateGeoNodesOptions struct {
	ID                               *int64    `url:"primary,omitempty" json:"primary,omitempty"`
	Enabled                          *bool     `url:"enabled,omitempty" json:"enabled,omitempty"`
	Name                             *string   `url:"name,omitempty" json:"name,omitempty"`
	URL                              *string   `url:"url,omitempty" json:"url,omitempty"`
	InternalURL                      *string   `url:"internal_url,omitempty" json:"internal_url,omitempty"`
	FilesMaxCapacity                 *int64    `url:"files_max_capacity,omitempty" json:"files_max_capacity,omitempty"`
	ReposMaxCapacity                 *int64    `url:"repos_max_capacity,omitempty" json:"repos_max_capacity,omitempty"`
	VerificationMaxCapacity          *int64    `url:"verification_max_capacity,omitempty" json:"verification_max_capacity,omitempty"`
	ContainerRepositoriesMaxCapacity *int64    `url:"container_repositories_max_capacity,omitempty" json:"container_repositories_max_capacity,omitempty"`
	SyncObjectStorage                *bool     `url:"sync_object_storage,omitempty" json:"sync_object_storage,omitempty"`
	SelectiveSyncType                *string   `url:"selective_sync_type,omitempty" json:"selective_sync_type,omitempty"`
	SelectiveSyncShards              *[]string `url:"selective_sync_shards,omitempty" json:"selective_sync_shards,omitempty"`
	SelectiveSyncNamespaceIDs        *[]int64  `url:"selective_sync_namespace_ids,omitempty" json:"selective_sync_namespace_ids,omitempty"`
	MinimumReverificationInterval    *int64    `url:"minimum_reverification_interval,omitempty" json:"minimum_reverification_interval,omitempty"`
}

// EditGeoNode updates settings of an existing Geo node.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#edit-a-geo-node
func (s *GeoNodesService) EditGeoNode(id int64, opt *UpdateGeoNodesOptions, options ...RequestOptionFunc) (*GeoNode, *Response, error) {
	return do[*GeoNode](s.client,
		withMethod(http.MethodPut),
		withPath("geo_nodes/%d", id),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteGeoNode removes the Geo node.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#delete-a-geo-node
func (s *GeoNodesService) DeleteGeoNode(id int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("geo_nodes/%d", id),
		withRequestOpts(options...),
	)
	return resp, err
}

// RepairGeoNode to repair the OAuth authentication of a Geo node.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#repair-a-geo-node
func (s *GeoNodesService) RepairGeoNode(id int64, options ...RequestOptionFunc) (*GeoNode, *Response, error) {
	return do[*GeoNode](s.client,
		withMethod(http.MethodPost),
		withPath("geo_nodes/%d/repair", id),
		withRequestOpts(options...),
	)
}

// GeoNodeStatus represents the status of Geo Node.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#retrieve-status-about-all-geo-nodes
type GeoNodeStatus struct {
	GeoNodeID                                     int64  `json:"geo_node_id"`
	Healthy                                       bool   `json:"healthy"`
	Health                                        string `json:"health"`
	HealthStatus                                  string `json:"health_status"`
	MissingOauthApplication                       bool   `json:"missing_oauth_application"`
	AttachmentsCount                              int64  `json:"attachments_count"`
	AttachmentsSyncedCount                        int64  `json:"attachments_synced_count"`
	AttachmentsFailedCount                        int64  `json:"attachments_failed_count"`
	AttachmentsSyncedMissingOnPrimaryCount        int64  `json:"attachments_synced_missing_on_primary_count"`
	AttachmentsSyncedInPercentage                 string `json:"attachments_synced_in_percentage"`
	DbReplicationLagSeconds                       int64  `json:"db_replication_lag_seconds"`
	LfsObjectsCount                               int64  `json:"lfs_objects_count"`
	LfsObjectsSyncedCount                         int64  `json:"lfs_objects_synced_count"`
	LfsObjectsFailedCount                         int64  `json:"lfs_objects_failed_count"`
	LfsObjectsSyncedMissingOnPrimaryCount         int64  `json:"lfs_objects_synced_missing_on_primary_count"`
	LfsObjectsSyncedInPercentage                  string `json:"lfs_objects_synced_in_percentage"`
	JobArtifactsCount                             int64  `json:"job_artifacts_count"`
	JobArtifactsSyncedCount                       int64  `json:"job_artifacts_synced_count"`
	JobArtifactsFailedCount                       int64  `json:"job_artifacts_failed_count"`
	JobArtifactsSyncedMissingOnPrimaryCount       int64  `json:"job_artifacts_synced_missing_on_primary_count"`
	JobArtifactsSyncedInPercentage                string `json:"job_artifacts_synced_in_percentage"`
	ContainerRepositoriesCount                    int64  `json:"container_repositories_count"`
	ContainerRepositoriesSyncedCount              int64  `json:"container_repositories_synced_count"`
	ContainerRepositoriesFailedCount              int64  `json:"container_repositories_failed_count"`
	ContainerRepositoriesSyncedInPercentage       string `json:"container_repositories_synced_in_percentage"`
	DesignRepositoriesCount                       int64  `json:"design_repositories_count"`
	DesignRepositoriesSyncedCount                 int64  `json:"design_repositories_synced_count"`
	DesignRepositoriesFailedCount                 int64  `json:"design_repositories_failed_count"`
	DesignRepositoriesSyncedInPercentage          string `json:"design_repositories_synced_in_percentage"`
	ProjectsCount                                 int64  `json:"projects_count"`
	RepositoriesCount                             int64  `json:"repositories_count"`
	RepositoriesFailedCount                       int64  `json:"repositories_failed_count"`
	RepositoriesSyncedCount                       int64  `json:"repositories_synced_count"`
	RepositoriesSyncedInPercentage                string `json:"repositories_synced_in_percentage"`
	WikisCount                                    int64  `json:"wikis_count"`
	WikisFailedCount                              int64  `json:"wikis_failed_count"`
	WikisSyncedCount                              int64  `json:"wikis_synced_count"`
	WikisSyncedInPercentage                       string `json:"wikis_synced_in_percentage"`
	ReplicationSlotsCount                         int64  `json:"replication_slots_count"`
	ReplicationSlotsUsedCount                     int64  `json:"replication_slots_used_count"`
	ReplicationSlotsUsedInPercentage              string `json:"replication_slots_used_in_percentage"`
	ReplicationSlotsMaxRetainedWalBytes           int64  `json:"replication_slots_max_retained_wal_bytes"`
	RepositoriesCheckedCount                      int64  `json:"repositories_checked_count"`
	RepositoriesCheckedFailedCount                int64  `json:"repositories_checked_failed_count"`
	RepositoriesCheckedInPercentage               string `json:"repositories_checked_in_percentage"`
	RepositoriesChecksummedCount                  int64  `json:"repositories_checksummed_count"`
	RepositoriesChecksumFailedCount               int64  `json:"repositories_checksum_failed_count"`
	RepositoriesChecksummedInPercentage           string `json:"repositories_checksummed_in_percentage"`
	WikisChecksummedCount                         int64  `json:"wikis_checksummed_count"`
	WikisChecksumFailedCount                      int64  `json:"wikis_checksum_failed_count"`
	WikisChecksummedInPercentage                  string `json:"wikis_checksummed_in_percentage"`
	RepositoriesVerifiedCount                     int64  `json:"repositories_verified_count"`
	RepositoriesVerificationFailedCount           int64  `json:"repositories_verification_failed_count"`
	RepositoriesVerifiedInPercentage              string `json:"repositories_verified_in_percentage"`
	RepositoriesChecksumMismatchCount             int64  `json:"repositories_checksum_mismatch_count"`
	WikisVerifiedCount                            int64  `json:"wikis_verified_count"`
	WikisVerificationFailedCount                  int64  `json:"wikis_verification_failed_count"`
	WikisVerifiedInPercentage                     string `json:"wikis_verified_in_percentage"`
	WikisChecksumMismatchCount                    int64  `json:"wikis_checksum_mismatch_count"`
	RepositoriesRetryingVerificationCount         int64  `json:"repositories_retrying_verification_count"`
	WikisRetryingVerificationCount                int64  `json:"wikis_retrying_verification_count"`
	LastEventID                                   int64  `json:"last_event_id"`
	LastEventTimestamp                            int64  `json:"last_event_timestamp"`
	CursorLastEventID                             int64  `json:"cursor_last_event_id"`
	CursorLastEventTimestamp                      int64  `json:"cursor_last_event_timestamp"`
	LastSuccessfulStatusCheckTimestamp            int64  `json:"last_successful_status_check_timestamp"`
	Version                                       string `json:"version"`
	Revision                                      string `json:"revision"`
	MergeRequestDiffsCount                        int64  `json:"merge_request_diffs_count"`
	MergeRequestDiffsChecksumTotalCount           int64  `json:"merge_request_diffs_checksum_total_count"`
	MergeRequestDiffsChecksummedCount             int64  `json:"merge_request_diffs_checksummed_count"`
	MergeRequestDiffsChecksumFailedCount          int64  `json:"merge_request_diffs_checksum_failed_count"`
	MergeRequestDiffsSyncedCount                  int64  `json:"merge_request_diffs_synced_count"`
	MergeRequestDiffsFailedCount                  int64  `json:"merge_request_diffs_failed_count"`
	MergeRequestDiffsRegistryCount                int64  `json:"merge_request_diffs_registry_count"`
	MergeRequestDiffsVerificationTotalCount       int64  `json:"merge_request_diffs_verification_total_count"`
	MergeRequestDiffsVerifiedCount                int64  `json:"merge_request_diffs_verified_count"`
	MergeRequestDiffsVerificationFailedCount      int64  `json:"merge_request_diffs_verification_failed_count"`
	MergeRequestDiffsSyncedInPercentage           string `json:"merge_request_diffs_synced_in_percentage"`
	MergeRequestDiffsVerifiedInPercentage         string `json:"merge_request_diffs_verified_in_percentage"`
	PackageFilesCount                             int64  `json:"package_files_count"`
	PackageFilesChecksumTotalCount                int64  `json:"package_files_checksum_total_count"`
	PackageFilesChecksummedCount                  int64  `json:"package_files_checksummed_count"`
	PackageFilesChecksumFailedCount               int64  `json:"package_files_checksum_failed_count"`
	PackageFilesSyncedCount                       int64  `json:"package_files_synced_count"`
	PackageFilesFailedCount                       int64  `json:"package_files_failed_count"`
	PackageFilesRegistryCount                     int64  `json:"package_files_registry_count"`
	PackageFilesVerificationTotalCount            int64  `json:"package_files_verification_total_count"`
	PackageFilesVerifiedCount                     int64  `json:"package_files_verified_count"`
	PackageFilesVerificationFailedCount           int64  `json:"package_files_verification_failed_count"`
	PackageFilesSyncedInPercentage                string `json:"package_files_synced_in_percentage"`
	PackageFilesVerifiedInPercentage              string `json:"package_files_verified_in_percentage"`
	PagesDeploymentsCount                         int64  `json:"pages_deployments_count"`
	PagesDeploymentsChecksumTotalCount            int64  `json:"pages_deployments_checksum_total_count"`
	PagesDeploymentsChecksummedCount              int64  `json:"pages_deployments_checksummed_count"`
	PagesDeploymentsChecksumFailedCount           int64  `json:"pages_deployments_checksum_failed_count"`
	PagesDeploymentsSyncedCount                   int64  `json:"pages_deployments_synced_count"`
	PagesDeploymentsFailedCount                   int64  `json:"pages_deployments_failed_count"`
	PagesDeploymentsRegistryCount                 int64  `json:"pages_deployments_registry_count"`
	PagesDeploymentsVerificationTotalCount        int64  `json:"pages_deployments_verification_total_count"`
	PagesDeploymentsVerifiedCount                 int64  `json:"pages_deployments_verified_count"`
	PagesDeploymentsVerificationFailedCount       int64  `json:"pages_deployments_verification_failed_count"`
	PagesDeploymentsSyncedInPercentage            string `json:"pages_deployments_synced_in_percentage"`
	PagesDeploymentsVerifiedInPercentage          string `json:"pages_deployments_verified_in_percentage"`
	TerraformStateVersionsCount                   int64  `json:"terraform_state_versions_count"`
	TerraformStateVersionsChecksumTotalCount      int64  `json:"terraform_state_versions_checksum_total_count"`
	TerraformStateVersionsChecksummedCount        int64  `json:"terraform_state_versions_checksummed_count"`
	TerraformStateVersionsChecksumFailedCount     int64  `json:"terraform_state_versions_checksum_failed_count"`
	TerraformStateVersionsSyncedCount             int64  `json:"terraform_state_versions_synced_count"`
	TerraformStateVersionsFailedCount             int64  `json:"terraform_state_versions_failed_count"`
	TerraformStateVersionsRegistryCount           int64  `json:"terraform_state_versions_registry_count"`
	TerraformStateVersionsVerificationTotalCount  int64  `json:"terraform_state_versions_verification_total_count"`
	TerraformStateVersionsVerifiedCount           int64  `json:"terraform_state_versions_verified_count"`
	TerraformStateVersionsVerificationFailedCount int64  `json:"terraform_state_versions_verification_failed_count"`
	TerraformStateVersionsSyncedInPercentage      string `json:"terraform_state_versions_synced_in_percentage"`
	TerraformStateVersionsVerifiedInPercentage    string `json:"terraform_state_versions_verified_in_percentage"`
	SnippetRepositoriesCount                      int64  `json:"snippet_repositories_count"`
	SnippetRepositoriesChecksumTotalCount         int64  `json:"snippet_repositories_checksum_total_count"`
	SnippetRepositoriesChecksummedCount           int64  `json:"snippet_repositories_checksummed_count"`
	SnippetRepositoriesChecksumFailedCount        int64  `json:"snippet_repositories_checksum_failed_count"`
	SnippetRepositoriesSyncedCount                int64  `json:"snippet_repositories_synced_count"`
	SnippetRepositoriesFailedCount                int64  `json:"snippet_repositories_failed_count"`
	SnippetRepositoriesRegistryCount              int64  `json:"snippet_repositories_registry_count"`
	SnippetRepositoriesVerificationTotalCount     int64  `json:"snippet_repositories_verification_total_count"`
	SnippetRepositoriesVerifiedCount              int64  `json:"snippet_repositories_verified_count"`
	SnippetRepositoriesVerificationFailedCount    int64  `json:"snippet_repositories_verification_failed_count"`
	SnippetRepositoriesSyncedInPercentage         string `json:"snippet_repositories_synced_in_percentage"`
	SnippetRepositoriesVerifiedInPercentage       string `json:"snippet_repositories_verified_in_percentage"`
	GroupWikiRepositoriesCount                    int64  `json:"group_wiki_repositories_count"`
	GroupWikiRepositoriesChecksumTotalCount       int64  `json:"group_wiki_repositories_checksum_total_count"`
	GroupWikiRepositoriesChecksummedCount         int64  `json:"group_wiki_repositories_checksummed_count"`
	GroupWikiRepositoriesChecksumFailedCount      int64  `json:"group_wiki_repositories_checksum_failed_count"`
	GroupWikiRepositoriesSyncedCount              int64  `json:"group_wiki_repositories_synced_count"`
	GroupWikiRepositoriesFailedCount              int64  `json:"group_wiki_repositories_failed_count"`
	GroupWikiRepositoriesRegistryCount            int64  `json:"group_wiki_repositories_registry_count"`
	GroupWikiRepositoriesVerificationTotalCount   int64  `json:"group_wiki_repositories_verification_total_count"`
	GroupWikiRepositoriesVerifiedCount            int64  `json:"group_wiki_repositories_verified_count"`
	GroupWikiRepositoriesVerificationFailedCount  int64  `json:"group_wiki_repositories_verification_failed_count"`
	GroupWikiRepositoriesSyncedInPercentage       string `json:"group_wiki_repositories_synced_in_percentage"`
	GroupWikiRepositoriesVerifiedInPercentage     string `json:"group_wiki_repositories_verified_in_percentage"`
	PipelineArtifactsCount                        int64  `json:"pipeline_artifacts_count"`
	PipelineArtifactsChecksumTotalCount           int64  `json:"pipeline_artifacts_checksum_total_count"`
	PipelineArtifactsChecksummedCount             int64  `json:"pipeline_artifacts_checksummed_count"`
	PipelineArtifactsChecksumFailedCount          int64  `json:"pipeline_artifacts_checksum_failed_count"`
	PipelineArtifactsSyncedCount                  int64  `json:"pipeline_artifacts_synced_count"`
	PipelineArtifactsFailedCount                  int64  `json:"pipeline_artifacts_failed_count"`
	PipelineArtifactsRegistryCount                int64  `json:"pipeline_artifacts_registry_count"`
	PipelineArtifactsVerificationTotalCount       int64  `json:"pipeline_artifacts_verification_total_count"`
	PipelineArtifactsVerifiedCount                int64  `json:"pipeline_artifacts_verified_count"`
	PipelineArtifactsVerificationFailedCount      int64  `json:"pipeline_artifacts_verification_failed_count"`
	PipelineArtifactsSyncedInPercentage           string `json:"pipeline_artifacts_synced_in_percentage"`
	PipelineArtifactsVerifiedInPercentage         string `json:"pipeline_artifacts_verified_in_percentage"`
	UploadsCount                                  int64  `json:"uploads_count"`
	UploadsSyncedCount                            int64  `json:"uploads_synced_count"`
	UploadsFailedCount                            int64  `json:"uploads_failed_count"`
	UploadsRegistryCount                          int64  `json:"uploads_registry_count"`
	UploadsSyncedInPercentage                     string `json:"uploads_synced_in_percentage"`
}

// RetrieveStatusOfAllGeoNodes get the list of status of all Geo Nodes.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#retrieve-status-about-all-geo-nodes
func (s *GeoNodesService) RetrieveStatusOfAllGeoNodes(options ...RequestOptionFunc) ([]*GeoNodeStatus, *Response, error) {
	return do[[]*GeoNodeStatus](s.client,
		withPath("geo_nodes/status"),
		withRequestOpts(options...),
	)
}

// RetrieveStatusOfGeoNode get the of status of a specific Geo Nodes.
// Deprecated: will be removed in v5 of the API, use Geo Sites API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_nodes/#retrieve-status-about-a-specific-geo-node
func (s *GeoNodesService) RetrieveStatusOfGeoNode(id int64, options ...RequestOptionFunc) (*GeoNodeStatus, *Response, error) {
	return do[*GeoNodeStatus](s.client,
		withPath("geo_nodes/%d/status", id),
		withRequestOpts(options...),
	)
}
