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

import (
	"net/http"
	"time"
)

type (
	GeoSitesServiceInterface interface {
		// CreateGeoSite creates a new Geo Site.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/geo_sites/#create-a-new-geo-site
		CreateGeoSite(*CreateGeoSitesOptions, ...RequestOptionFunc) (*GeoSite, *Response, error)
		// ListGeoSites gets a list of geo sites.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/geo_sites/#retrieve-configuration-about-all-geo-sites
		ListGeoSites(*ListGeoSitesOptions, ...RequestOptionFunc) ([]*GeoSite, *Response, error)
		// GetGeoSite gets a specific geo site.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/geo_sites/#retrieve-configuration-about-a-specific-geo-site
		GetGeoSite(int64, ...RequestOptionFunc) (*GeoSite, *Response, error)
		// EditGeoSite updates settings of an existing Geo site.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/geo_sites/#edit-a-geo-site
		EditGeoSite(int64, *EditGeoSiteOptions, ...RequestOptionFunc) (*GeoSite, *Response, error)
		// DeleteGeoSite removes the Geo site.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/geo_sites/#delete-a-geo-site
		DeleteGeoSite(int64, ...RequestOptionFunc) (*Response, error)
		// RepairGeoSite to repair the OAuth authentication of a Geo site.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/geo_sites/#repair-a-geo-site
		RepairGeoSite(int64, ...RequestOptionFunc) (*GeoSite, *Response, error)
		// ListStatusOfAllGeoSites get the list of status of all Geo Sites.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/geo_sites/#retrieve-status-about-all-geo-sites
		ListStatusOfAllGeoSites(*ListStatusOfAllGeoSitesOptions, ...RequestOptionFunc) ([]*GeoSiteStatus, *Response, error)
		// GetStatusOfGeoSite gets the status of a specific Geo Site.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/geo_sites/#retrieve-status-about-a-specific-geo-site
		GetStatusOfGeoSite(int64, ...RequestOptionFunc) (*GeoSiteStatus, *Response, error)
	}

	// GeoSitesService handles communication with Geo Sites related methods
	// of GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/geo_sites/
	GeoSitesService struct {
		client *Client
	}
)

var _ GeoSitesServiceInterface = (*GeoSitesService)(nil)

// GeoSite represents a GitLab Geo Site.
//
// GitLab API docs: https://docs.gitlab.com/api/geo_sites/
type GeoSite struct {
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
	ContainerRepositoriesMaxCapacity int64        `json:"container_repositories_max_capacity"`
	SelectiveSyncType                string       `json:"selective_sync_type"`
	SelectiveSyncShards              []string     `json:"selective_sync_shards"`
	SelectiveSyncNamespaceIDs        []int64      `json:"selective_sync_namespace_ids"`
	MinimumReverificationInterval    int64        `json:"minimum_reverification_interval"`
	SyncObjectStorage                bool         `json:"sync_object_storage"`
	WebEditURL                       string       `json:"web_edit_url"`
	WebGeoReplicationDetailsURL      string       `json:"web_geo_replication_details_url"`
	Links                            GeoSiteLinks `json:"_links"`
}

// GeoSiteLinks represents links for GitLab GeoSite.
//
// GitLab API docs: https://docs.gitlab.com/api/geo_sites/
type GeoSiteLinks struct {
	Self   string `json:"self"`
	Status string `json:"status"`
	Repair string `json:"repair"`
}

// CreateGeoSitesOptions represents the available CreateGeoSite() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#create-a-geo-site
type CreateGeoSitesOptions struct {
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

func (s *GeoSitesService) CreateGeoSite(opt *CreateGeoSitesOptions, options ...RequestOptionFunc) (*GeoSite, *Response, error) {
	return do[*GeoSite](s.client,
		withMethod(http.MethodPost),
		withPath("geo_sites"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListGeoSitesOptions represents the available ListGeoSites() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#list-all-geo-sites
type ListGeoSitesOptions struct {
	ListOptions
}

func (s *GeoSitesService) ListGeoSites(opt *ListGeoSitesOptions, options ...RequestOptionFunc) ([]*GeoSite, *Response, error) {
	return do[[]*GeoSite](s.client,
		withPath("geo_sites"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *GeoSitesService) GetGeoSite(id int64, options ...RequestOptionFunc) (*GeoSite, *Response, error) {
	return do[*GeoSite](s.client,
		withPath("geo_sites/%d", id),
		withRequestOpts(options...),
	)
}

// EditGeoSiteOptions represents the available EditGeoSite() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#update-a-geo-site
type EditGeoSiteOptions struct {
	Enabled                          *bool     `url:"enabled,omitempty" json:"enabled,omitempty"`
	Name                             *string   `url:"name,omitempty" json:"name,omitempty"`
	URL                              *string   `url:"url,omitempty" json:"url,omitempty"`
	InternalURL                      *string   `url:"internal_url,omitempty" json:"internal_url,omitempty"`
	FilesMaxCapacity                 *int64    `url:"files_max_capacity,omitempty" json:"files_max_capacity,omitempty"`
	ReposMaxCapacity                 *int64    `url:"repos_max_capacity,omitempty" json:"repos_max_capacity,omitempty"`
	VerificationMaxCapacity          *int64    `url:"verification_max_capacity,omitempty" json:"verification_max_capacity,omitempty"`
	ContainerRepositoriesMaxCapacity *int64    `url:"container_repositories_max_capacity,omitempty" json:"container_repositories_max_capacity,omitempty"`
	SelectiveSyncType                *string   `url:"selective_sync_type,omitempty" json:"selective_sync_type,omitempty"`
	SelectiveSyncShards              *[]string `url:"selective_sync_shards,omitempty" json:"selective_sync_shards,omitempty"`
	SelectiveSyncNamespaceIDs        *[]int64  `url:"selective_sync_namespace_ids,omitempty" json:"selective_sync_namespace_ids,omitempty"`
	MinimumReverificationInterval    *int64    `url:"minimum_reverification_interval,omitempty" json:"minimum_reverification_interval,omitempty"`
}

func (s *GeoSitesService) EditGeoSite(id int64, opt *EditGeoSiteOptions, options ...RequestOptionFunc) (*GeoSite, *Response, error) {
	return do[*GeoSite](s.client,
		withMethod(http.MethodPut),
		withPath("geo_sites/%d", id),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *GeoSitesService) DeleteGeoSite(id int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("geo_sites/%d", id),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *GeoSitesService) RepairGeoSite(id int64, options ...RequestOptionFunc) (*GeoSite, *Response, error) {
	return do[*GeoSite](s.client,
		withMethod(http.MethodPost),
		withPath("geo_sites/%d/repair", id),
		withRequestOpts(options...),
	)
}

// GeoSiteStatus represents the status of Geo Site.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#list-all-geo-site-statuses
type GeoSiteStatus struct {
	GeoNodeID                                       int64             `json:"geo_node_id"`
	ProjectsCount                                   int64             `json:"projects_count"`
	ContainerRepositoriesReplicationEnabled         bool              `json:"container_repositories_replication_enabled"`
	LFSObjectsCount                                 int64             `json:"lfs_objects_count"`
	LFSObjectsChecksumTotalCount                    int64             `json:"lfs_objects_checksum_total_count"`
	LFSObjectsChecksummedCount                      int64             `json:"lfs_objects_checksummed_count"`
	LFSObjectsChecksumFailedCount                   int64             `json:"lfs_objects_checksum_failed_count"`
	LFSObjectsSyncedCount                           int64             `json:"lfs_objects_synced_count"`
	LFSObjectsFailedCount                           int64             `json:"lfs_objects_failed_count"`
	LFSObjectsRegistryCount                         int64             `json:"lfs_objects_registry_count"`
	LFSObjectsVerificationTotalCount                int64             `json:"lfs_objects_verification_total_count"`
	LFSObjectsVerifiedCount                         int64             `json:"lfs_objects_verified_count"`
	LFSObjectsVerificationFailedCount               int64             `json:"lfs_objects_verification_failed_count"`
	MergeRequestDiffsCount                          int64             `json:"merge_request_diffs_count"`
	MergeRequestDiffsChecksumTotalCount             int64             `json:"merge_request_diffs_checksum_total_count"`
	MergeRequestDiffsChecksummedCount               int64             `json:"merge_request_diffs_checksummed_count"`
	MergeRequestDiffsChecksumFailedCount            int64             `json:"merge_request_diffs_checksum_failed_count"`
	MergeRequestDiffsSyncedCount                    int64             `json:"merge_request_diffs_synced_count"`
	MergeRequestDiffsFailedCount                    int64             `json:"merge_request_diffs_failed_count"`
	MergeRequestDiffsRegistryCount                  int64             `json:"merge_request_diffs_registry_count"`
	MergeRequestDiffsVerificationTotalCount         int64             `json:"merge_request_diffs_verification_total_count"`
	MergeRequestDiffsVerifiedCount                  int64             `json:"merge_request_diffs_verified_count"`
	MergeRequestDiffsVerificationFailedCount        int64             `json:"merge_request_diffs_verification_failed_count"`
	PackageFilesCount                               int64             `json:"package_files_count"`
	PackageFilesChecksumTotalCount                  int64             `json:"package_files_checksum_total_count"`
	PackageFilesChecksummedCount                    int64             `json:"package_files_checksummed_count"`
	PackageFilesChecksumFailedCount                 int64             `json:"package_files_checksum_failed_count"`
	PackageFilesSyncedCount                         int64             `json:"package_files_synced_count"`
	PackageFilesFailedCount                         int64             `json:"package_files_failed_count"`
	PackageFilesRegistryCount                       int64             `json:"package_files_registry_count"`
	PackageFilesVerificationTotalCount              int64             `json:"package_files_verification_total_count"`
	PackageFilesVerifiedCount                       int64             `json:"package_files_verified_count"`
	PackageFilesVerificationFailedCount             int64             `json:"package_files_verification_failed_count"`
	TerraformStateVersionsCount                     int64             `json:"terraform_state_versions_count"`
	TerraformStateVersionsChecksumTotalCount        int64             `json:"terraform_state_versions_checksum_total_count"`
	TerraformStateVersionsChecksummedCount          int64             `json:"terraform_state_versions_checksummed_count"`
	TerraformStateVersionsChecksumFailedCount       int64             `json:"terraform_state_versions_checksum_failed_count"`
	TerraformStateVersionsSyncedCount               int64             `json:"terraform_state_versions_synced_count"`
	TerraformStateVersionsFailedCount               int64             `json:"terraform_state_versions_failed_count"`
	TerraformStateVersionsRegistryCount             int64             `json:"terraform_state_versions_registry_count"`
	TerraformStateVersionsVerificationTotalCount    int64             `json:"terraform_state_versions_verification_total_count"`
	TerraformStateVersionsVerifiedCount             int64             `json:"terraform_state_versions_verified_count"`
	TerraformStateVersionsVerificationFailedCount   int64             `json:"terraform_state_versions_verification_failed_count"`
	SnippetRepositoriesCount                        int64             `json:"snippet_repositories_count"`
	SnippetRepositoriesChecksumTotalCount           int64             `json:"snippet_repositories_checksum_total_count"`
	SnippetRepositoriesChecksummedCount             int64             `json:"snippet_repositories_checksummed_count"`
	SnippetRepositoriesChecksumFailedCount          int64             `json:"snippet_repositories_checksum_failed_count"`
	SnippetRepositoriesSyncedCount                  int64             `json:"snippet_repositories_synced_count"`
	SnippetRepositoriesFailedCount                  int64             `json:"snippet_repositories_failed_count"`
	SnippetRepositoriesRegistryCount                int64             `json:"snippet_repositories_registry_count"`
	SnippetRepositoriesVerificationTotalCount       int64             `json:"snippet_repositories_verification_total_count"`
	SnippetRepositoriesVerifiedCount                int64             `json:"snippet_repositories_verified_count"`
	SnippetRepositoriesVerificationFailedCount      int64             `json:"snippet_repositories_verification_failed_count"`
	GroupWikiRepositoriesCount                      int64             `json:"group_wiki_repositories_count"`
	GroupWikiRepositoriesChecksumTotalCount         int64             `json:"group_wiki_repositories_checksum_total_count"`
	GroupWikiRepositoriesChecksummedCount           int64             `json:"group_wiki_repositories_checksummed_count"`
	GroupWikiRepositoriesChecksumFailedCount        int64             `json:"group_wiki_repositories_checksum_failed_count"`
	GroupWikiRepositoriesSyncedCount                int64             `json:"group_wiki_repositories_synced_count"`
	GroupWikiRepositoriesFailedCount                int64             `json:"group_wiki_repositories_failed_count"`
	GroupWikiRepositoriesRegistryCount              int64             `json:"group_wiki_repositories_registry_count"`
	GrupWikiRepositoriesVerificationTotalCount      int64             `json:"group_wiki_repositories_verification_total_count"`
	GroupWikiRepositoriesVerifiedCount              int64             `json:"group_wiki_repositories_verified_count"`
	GroupWikiRepositoriesVerificationFailedCount    int64             `json:"group_wiki_repositories_verification_failed_count"`
	PipelineArtifactsCount                          int64             `json:"pipeline_artifacts_count"`
	PipelineArtifactsChecksumTotalCount             int64             `json:"pipeline_artifacts_checksum_total_count"`
	PipelineArtifactsChecksummedCount               int64             `json:"pipeline_artifacts_checksummed_count"`
	PipelineArtifactsChecksumFailedCount            int64             `json:"pipeline_artifacts_checksum_failed_count"`
	PipelineArtifactsSyncedCount                    int64             `json:"pipeline_artifacts_synced_count"`
	PipelineArtifactsFailedCount                    int64             `json:"pipeline_artifacts_failed_count"`
	PipelineArtifactsRegistryCount                  int64             `json:"pipeline_artifacts_registry_count"`
	PipelineArtifactsVerificationTotalCount         int64             `json:"pipeline_artifacts_verification_total_count"`
	PipelineArtifactsVerifiedCount                  int64             `json:"pipeline_artifacts_verified_count"`
	PipelineArtifactsVerificationFailedCount        int64             `json:"pipeline_artifacts_verification_failed_count"`
	PagesDeploymentsCount                           int64             `json:"pages_deployments_count"`
	PagesDeploymentsChecksumTotalCount              int64             `json:"pages_deployments_checksum_total_count"`
	PagesDeploymentsChecksummedCount                int64             `json:"pages_deployments_checksummed_count"`
	PagesDeploymentsChecksumFailedCount             int64             `json:"pages_deployments_checksum_failed_count"`
	PagesDeploymentsSyncedCount                     int64             `json:"pages_deployments_synced_count"`
	PagesDeploymentsFailedCount                     int64             `json:"pages_deployments_failed_count"`
	PagesDeploymentsRegistryCount                   int64             `json:"pages_deployments_registry_count"`
	PagesDeploymentsVerificationTotalCount          int64             `json:"pages_deployments_verification_total_count"`
	PagesDeploymentsVerifiedCount                   int64             `json:"pages_deployments_verified_count"`
	PagesDeploymentsVerificationFailedCount         int64             `json:"pages_deployments_verification_failed_count"`
	UploadsCount                                    int64             `json:"uploads_count"`
	UploadsChecksumTotalCount                       int64             `json:"uploads_checksum_total_count"`
	UploadsChecksummedCount                         int64             `json:"uploads_checksummed_count"`
	UploadsChecksumFailedCount                      int64             `json:"uploads_checksum_failed_count"`
	UploadsSyncedCount                              int64             `json:"uploads_synced_count"`
	UploadsFailedCount                              int64             `json:"uploads_failed_count"`
	UploadsRegistryCount                            int64             `json:"uploads_registry_count"`
	UploadsVerificationTotalCount                   int64             `json:"uploads_verification_total_count"`
	UploadsVerifiedCount                            int64             `json:"uploads_verified_count"`
	UploadsVerificationFailedCount                  int64             `json:"uploads_verification_failed_count"`
	JobArtifactsCount                               int64             `json:"job_artifacts_count"`
	JobArtifactsChecksumTotalCount                  int64             `json:"job_artifacts_checksum_total_count"`
	JobArtifactsChecksummedCount                    int64             `json:"job_artifacts_checksummed_count"`
	JobArtifactsChecksumFailedCount                 int64             `json:"job_artifacts_checksum_failed_count"`
	JobArtifactsSyncedCount                         int64             `json:"job_artifacts_synced_count"`
	JobArtifactsFailedCount                         int64             `json:"job_artifacts_failed_count"`
	JobArtifactsRegistryCount                       int64             `json:"job_artifacts_registry_count"`
	JobArtifactsVerificationTotalCount              int64             `json:"job_artifacts_verification_total_count"`
	JobArtifactsVerifiedCount                       int64             `json:"job_artifacts_verified_count"`
	JobArtifactsVerificationFailedCount             int64             `json:"job_artifacts_verification_failed_count"`
	CISecureFilesCount                              int64             `json:"ci_secure_files_count"`
	CISecureFilesChecksumTotalCount                 int64             `json:"ci_secure_files_checksum_total_count"`
	CISecureFilesChecksummedCount                   int64             `json:"ci_secure_files_checksummed_count"`
	CISecureFilesChecksumFailedCount                int64             `json:"ci_secure_files_checksum_failed_count"`
	CISecureFilesSyncedCount                        int64             `json:"ci_secure_files_synced_count"`
	CISecureFilesFailedCount                        int64             `json:"ci_secure_files_failed_count"`
	CISecureFilesRegistryCount                      int64             `json:"ci_secure_files_registry_count"`
	CISecureFilesVerificationTotalCount             int64             `json:"ci_secure_files_verification_total_count"`
	CISecureFilesVerifiedCount                      int64             `json:"ci_secure_files_verified_count"`
	CISecureFilesVerificationFailedCount            int64             `json:"ci_secure_files_verification_failed_count"`
	ContainerRepositoriesCount                      int64             `json:"container_repositories_count"`
	ContainerRepositoriesChecksumTotalCount         int64             `json:"container_repositories_checksum_total_count"`
	ContainerRepositoriesChecksummedCount           int64             `json:"container_repositories_checksummed_count"`
	ContainerRepositoriesChecksumFailedCount        int64             `json:"container_repositories_checksum_failed_count"`
	ContainerRepositoriesSyncedCount                int64             `json:"container_repositories_synced_count"`
	ContainerRepositoriesFailedCount                int64             `json:"container_repositories_failed_count"`
	ContainerRepositoriesRegistryCount              int64             `json:"container_repositories_registry_count"`
	ContainerRepositoriesVerificationTotalCount     int64             `json:"container_repositories_verification_total_count"`
	ContainerRepositoriesVerifiedCount              int64             `json:"container_repositories_verified_count"`
	ContainerRepositoriesVerificationFailedCount    int64             `json:"container_repositories_verification_failed_count"`
	DependencyProxyBlobsCount                       int64             `json:"dependency_proxy_blobs_count"`
	DependencyProxyBlobsChecksumTotalCount          int64             `json:"dependency_proxy_blobs_checksum_total_count"`
	DependencyProxyBlobsChecksummedCount            int64             `json:"dependency_proxy_blobs_checksummed_count"`
	DependencyProxyBlobsChecksumFailedCount         int64             `json:"dependency_proxy_blobs_checksum_failed_count"`
	DependencyProxyBlobsSyncedCount                 int64             `json:"dependency_proxy_blobs_synced_count"`
	DependencyProxyBlobsFailedCount                 int64             `json:"dependency_proxy_blobs_failed_count"`
	DependencyProxyBlobsRegistryCount               int64             `json:"dependency_proxy_blobs_registry_count"`
	DependencyProxyBlobsVerificationTotalCount      int64             `json:"dependency_proxy_blobs_verification_total_count"`
	DependencyProxyBlobsVerifiedCount               int64             `json:"dependency_proxy_blobs_verified_count"`
	DependencyProxyBlobsVerificationFailedCount     int64             `json:"dependency_proxy_blobs_verification_failed_count"`
	DependencyProxyManifestsCount                   int64             `json:"dependency_proxy_manifests_count"`
	DependencyProxyManifestsChecksumTotalCount      int64             `json:"dependency_proxy_manifests_checksum_total_count"`
	DependencyProxyManifestsChecksummedCount        int64             `json:"dependency_proxy_manifests_checksummed_count"`
	DependencyProxyManifestsChecksumFailedCount     int64             `json:"dependency_proxy_manifests_checksum_failed_count"`
	DependencyProxyManifestsSyncedCount             int64             `json:"dependency_proxy_manifests_synced_count"`
	DependencyProxyManifestsFailedCount             int64             `json:"dependency_proxy_manifests_failed_count"`
	DependencyProxyManifestsRegistryCount           int64             `json:"dependency_proxy_manifests_registry_count"`
	DependencyProxyManifestsVerificationTotalCount  int64             `json:"dependency_proxy_manifests_verification_total_count"`
	DependencyProxyManifestsVerifiedCount           int64             `json:"dependency_proxy_manifests_verified_count"`
	DependencyProxyManifestsVerificationFailedCount int64             `json:"dependency_proxy_manifests_verification_failed_count"`
	ProjectWikiRepositoriesCount                    int64             `json:"project_wiki_repositories_count"`
	ProjectWikiRepositoriesChecksumTotalCount       int64             `json:"project_wiki_repositories_checksum_total_count"`
	ProjectWikiRepositoriesChecksummedCount         int64             `json:"project_wiki_repositories_checksummed_count"`
	ProjectWikiRepositoriesChecksumFailedCount      int64             `json:"project_wiki_repositories_checksum_failed_count"`
	ProjectWikiRepositoriesSyncedCount              int64             `json:"project_wiki_repositories_synced_count"`
	ProjectWikiRepositoriesFailedCount              int64             `json:"project_wiki_repositories_failed_count"`
	ProjectWikiRepositoriesRegistryCount            int64             `json:"project_wiki_repositories_registry_count"`
	ProjectWikiRepositoriesVerificationTotalCount   int64             `json:"project_wiki_repositories_verification_total_count"`
	ProjectWikiRepositoriesVerifiedCount            int64             `json:"project_wiki_repositories_verified_count"`
	ProjectWikiRepositoriesVerificationFailedCount  int64             `json:"project_wiki_repositories_verification_failed_count"`
	GitFetchEventCountWeekly                        int64             `json:"git_fetch_event_count_weekly"`
	GitPushEventCountWeekly                         int64             `json:"git_push_event_count_weekly"`
	ProxyRemoteRequestsEventCountWeekly             int64             `json:"proxy_remote_requests_event_count_weekly"`
	ProxyLocalRequestsEventCountWeekly              int64             `json:"proxy_local_requests_event_count_weekly"`
	RepositoriesCheckedInPercentage                 string            `json:"repositories_checked_in_percentage"`
	ReplicationSlotsUsedInPercentage                string            `json:"replication_slots_used_in_percentage"`
	LFSObjectsSyncedInPercentage                    string            `json:"lfs_objects_synced_in_percentage"`
	LFSObjectsVerifiedInPercentage                  string            `json:"lfs_objects_verified_in_percentage"`
	MergeRequestDiffsSyncedInPercentage             string            `json:"merge_request_diffs_synced_in_percentage"`
	MergeRequestDiffsVerifiedInPercentage           string            `json:"merge_request_diffs_verified_in_percentage"`
	PackageFilesSyncedInPercentage                  string            `json:"package_files_synced_in_percentage"`
	PackageFilesVerifiedInPercentage                string            `json:"package_files_verified_in_percentage"`
	TerraformStateVersionsSyncedInPercentage        string            `json:"terraform_state_versions_synced_in_percentage"`
	TerraformStateVersionsVerifiedInPercentage      string            `json:"terraform_state_versions_verified_in_percentage"`
	SnippetRepositoriesSyncedInPercentage           string            `json:"snippet_repositories_synced_in_percentage"`
	SnippetRepositoriesVerifiedInPercentage         string            `json:"snippet_repositories_verified_in_percentage"`
	GroupWikiRepositoriesSyncedInPercentage         string            `json:"group_wiki_repositories_synced_in_percentage"`
	GroupWikiRepositoriesVerifiedInPercentage       string            `json:"group_wiki_repositories_verified_in_percentage"`
	PipelineArtifactsSyncedInPercentage             string            `json:"pipeline_artifacts_synced_in_percentage"`
	PipelineArtifactsVerifiedInPercentage           string            `json:"pipeline_artifacts_verified_in_percentage"`
	PagesDeploymentsSyncedInPercentage              string            `json:"pages_deployments_synced_in_percentage"`
	PagesDeploymentsVerifiedInPercentage            string            `json:"pages_deployments_verified_in_percentage"`
	UploadsSyncedInPercentage                       string            `json:"uploads_synced_in_percentage"`
	UploadsVerifiedInPercentage                     string            `json:"uploads_verified_in_percentage"`
	JobArtifactsSyncedInPercentage                  string            `json:"job_artifacts_synced_in_percentage"`
	JobArtifactsVerifiedInPercentage                string            `json:"job_artifacts_verified_in_percentage"`
	CISecureFilesSyncedInPercentage                 string            `json:"ci_secure_files_synced_in_percentage"`
	CISecureFilesVerifiedInPercentage               string            `json:"ci_secure_files_verified_in_percentage"`
	ContainerRepositoriesSyncedInPercentage         string            `json:"container_repositories_synced_in_percentage"`
	ContainerRepositoriesVerifiedInPercentage       string            `json:"container_repositories_verified_in_percentage"`
	DependencyProxyBlobsSyncedInPercentage          string            `json:"dependency_proxy_blobs_synced_in_percentage"`
	DependencyProxyBlobsVerifiedInPercentage        string            `json:"dependency_proxy_blobs_verified_in_percentage"`
	DependencyProxyManifestsSyncedInPercentage      string            `json:"dependency_proxy_manifests_synced_in_percentage"`
	DependencyProxyManifestsVerifiedInPercentage    string            `json:"dependency_proxy_manifests_verified_in_percentage"`
	ProjectWikiRepositoriesSyncedInPercentage       string            `json:"project_wiki_repositories_synced_in_percentage"`
	ProjectWikiRepositoriesVerifiedInPercentage     string            `json:"project_wiki_repositories_verified_in_percentage"`
	ReplicationSlotsCount                           int64             `json:"replication_slots_count"`
	ReplicationSlotsUsedCount                       int64             `json:"replication_slots_used_count"`
	Healthy                                         bool              `json:"healthy"`
	Health                                          string            `json:"health"`
	HealthStatus                                    string            `json:"health_status"`
	MissingOAuthApplication                         bool              `json:"missing_oauth_application"`
	DBReplicationLagSeconds                         int64             `json:"db_replication_lag_seconds"`
	ReplicationSlotsMaxRetainedWalBytes             int64             `json:"replication_slots_max_retained_wal_bytes"`
	RepositoriesCheckedCount                        int64             `json:"repositories_checked_count"`
	RepositoriesCheckedFailedCount                  int64             `json:"repositories_checked_failed_count"`
	LastEventID                                     int64             `json:"last_event_id"`
	LastEventTimestamp                              int64             `json:"last_event_timestamp"`
	CursorLastEventID                               int64             `json:"cursor_last_event_id"`
	CursorLastEventTimestamp                        int64             `json:"cursor_last_event_timestamp"`
	LastSuccessfulStatusCheckTimestamp              int64             `json:"last_successful_status_check_timestamp"`
	Version                                         string            `json:"version"`
	Revision                                        string            `json:"revision"`
	SelectiveSyncType                               string            `json:"selective_sync_type"`
	Namespaces                                      []string          `json:"namespaces"`
	UpdatedAt                                       time.Time         `json:"updated_at"`
	StorageShardsMatch                              bool              `json:"storage_shards_match"`
	Links                                           GeoSiteStatusLink `json:"_links"`
}

// GeoSiteStatusLink represents the links for a GitLab Geo Site status.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#list-all-geo-site-statuses
type GeoSiteStatusLink struct {
	Self string `json:"self"`
	Site string `json:"site"`
}

// ListStatusOfAllGeoSitesOptions represents the available ListStatusOfAllGeoSites() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#list-all-geo-site-statuses
type ListStatusOfAllGeoSitesOptions struct {
	ListOptions
}

func (s *GeoSitesService) ListStatusOfAllGeoSites(opt *ListStatusOfAllGeoSitesOptions, options ...RequestOptionFunc) ([]*GeoSiteStatus, *Response, error) {
	return do[[]*GeoSiteStatus](s.client,
		withPath("geo_sites/status"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *GeoSitesService) GetStatusOfGeoSite(id int64, options ...RequestOptionFunc) (*GeoSiteStatus, *Response, error) {
	return do[*GeoSiteStatus](s.client,
		withPath("geo_sites/%d/status", id),
		withRequestOpts(options...),
	)
}
