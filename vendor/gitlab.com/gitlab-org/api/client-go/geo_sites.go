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
	"fmt"
	"net/http"
	"time"
)

type (
	GeoSitesServiceInterface interface {
		CreateGeoSite(*CreateGeoSitesOptions, ...RequestOptionFunc) (*GeoSite, *Response, error)
		ListGeoSites(*ListGeoSitesOptions, ...RequestOptionFunc) ([]*GeoSite, *Response, error)
		GetGeoSite(int, ...RequestOptionFunc) (*GeoSite, *Response, error)
		EditGeoSite(int, *EditGeoSiteOptions, ...RequestOptionFunc) (*GeoSite, *Response, error)
		DeleteGeoSite(int, ...RequestOptionFunc) (*Response, error)
		RepairGeoSite(int, ...RequestOptionFunc) (*GeoSite, *Response, error)
		ListStatusOfAllGeoSites(*ListStatusOfAllGeoSitesOptions, ...RequestOptionFunc) ([]*GeoSiteStatus, *Response, error)
		GetStatusOfGeoSite(int, ...RequestOptionFunc) (*GeoSiteStatus, *Response, error)
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
	ID                               int          `json:"id"`
	Name                             string       `json:"name"`
	URL                              string       `json:"url"`
	InternalURL                      string       `json:"internal_url"`
	Primary                          bool         `json:"primary"`
	Enabled                          bool         `json:"enabled"`
	Current                          bool         `json:"current"`
	FilesMaxCapacity                 int          `json:"files_max_capacity"`
	ReposMaxCapacity                 int          `json:"repos_max_capacity"`
	VerificationMaxCapacity          int          `json:"verification_max_capacity"`
	ContainerRepositoriesMaxCapacity int          `json:"container_repositories_max_capacity"`
	SelectiveSyncType                string       `json:"selective_sync_type"`
	SelectiveSyncShards              []string     `json:"selective_sync_shards"`
	SelectiveSyncNamespaceIDs        []int        `json:"selective_sync_namespace_ids"`
	MinimumReverificationInterval    int          `json:"minimum_reverification_interval"`
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
// https://docs.gitlab.com/api/geo_sites/#create-a-new-geo-site
type CreateGeoSitesOptions struct {
	Primary                          *bool     `url:"primary,omitempty" json:"primary,omitempty"`
	Enabled                          *bool     `url:"enabled,omitempty" json:"enabled,omitempty"`
	Name                             *string   `url:"name,omitempty" json:"name,omitempty"`
	URL                              *string   `url:"url,omitempty" json:"url,omitempty"`
	InternalURL                      *string   `url:"internal_url,omitempty" json:"internal_url,omitempty"`
	FilesMaxCapacity                 *int      `url:"files_max_capacity,omitempty" json:"files_max_capacity,omitempty"`
	ReposMaxCapacity                 *int      `url:"repos_max_capacity,omitempty" json:"repos_max_capacity,omitempty"`
	VerificationMaxCapacity          *int      `url:"verification_max_capacity,omitempty" json:"verification_max_capacity,omitempty"`
	ContainerRepositoriesMaxCapacity *int      `url:"container_repositories_max_capacity,omitempty" json:"container_repositories_max_capacity,omitempty"`
	SyncObjectStorage                *bool     `url:"sync_object_storage,omitempty" json:"sync_object_storage,omitempty"`
	SelectiveSyncType                *string   `url:"selective_sync_type,omitempty" json:"selective_sync_type,omitempty"`
	SelectiveSyncShards              *[]string `url:"selective_sync_shards,omitempty" json:"selective_sync_shards,omitempty"`
	SelectiveSyncNamespaceIDs        *[]int    `url:"selective_sync_namespace_ids,omitempty" json:"selective_sync_namespace_ids,omitempty"`
	MinimumReverificationInterval    *int      `url:"minimum_reverification_interval,omitempty" json:"minimum_reverification_interval,omitempty"`
}

// CreateGeoSite creates a new Geo Site.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#create-a-new-geo-site
func (s *GeoSitesService) CreateGeoSite(opt *CreateGeoSitesOptions, options ...RequestOptionFunc) (*GeoSite, *Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "geo_sites", opt, options)
	if err != nil {
		return nil, nil, err
	}

	site := new(GeoSite)
	resp, err := s.client.Do(req, site)
	if err != nil {
		return nil, resp, err
	}

	return site, resp, nil
}

// ListGeoSitesOptions represents the available ListGeoSites() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#retrieve-configuration-about-all-geo-sites
type ListGeoSitesOptions ListOptions

// ListGeoSites gets a list of geo sites.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#retrieve-configuration-about-all-geo-sites
func (s *GeoSitesService) ListGeoSites(opt *ListGeoSitesOptions, options ...RequestOptionFunc) ([]*GeoSite, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "geo_sites", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var sites []*GeoSite
	resp, err := s.client.Do(req, &sites)
	if err != nil {
		return nil, resp, err
	}

	return sites, resp, nil
}

// GetGeoSite gets a specific geo site.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#retrieve-configuration-about-a-specific-geo-site
func (s *GeoSitesService) GetGeoSite(id int, options ...RequestOptionFunc) (*GeoSite, *Response, error) {
	u := fmt.Sprintf("geo_sites/%d", id)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	site := new(GeoSite)
	resp, err := s.client.Do(req, site)
	if err != nil {
		return nil, resp, err
	}

	return site, resp, nil
}

// EditGeoSiteOptions represents the available EditGeoSite() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#edit-a-geo-site
type EditGeoSiteOptions struct {
	Enabled                          *bool     `url:"enabled,omitempty" json:"enabled,omitempty"`
	Name                             *string   `url:"name,omitempty" json:"name,omitempty"`
	URL                              *string   `url:"url,omitempty" json:"url,omitempty"`
	InternalURL                      *string   `url:"internal_url,omitempty" json:"internal_url,omitempty"`
	FilesMaxCapacity                 *int      `url:"files_max_capacity,omitempty" json:"files_max_capacity,omitempty"`
	ReposMaxCapacity                 *int      `url:"repos_max_capacity,omitempty" json:"repos_max_capacity,omitempty"`
	VerificationMaxCapacity          *int      `url:"verification_max_capacity,omitempty" json:"verification_max_capacity,omitempty"`
	ContainerRepositoriesMaxCapacity *int      `url:"container_repositories_max_capacity,omitempty" json:"container_repositories_max_capacity,omitempty"`
	SelectiveSyncType                *string   `url:"selective_sync_type,omitempty" json:"selective_sync_type,omitempty"`
	SelectiveSyncShards              *[]string `url:"selective_sync_shards,omitempty" json:"selective_sync_shards,omitempty"`
	SelectiveSyncNamespaceIDs        *[]int    `url:"selective_sync_namespace_ids,omitempty" json:"selective_sync_namespace_ids,omitempty"`
	MinimumReverificationInterval    *int      `url:"minimum_reverification_interval,omitempty" json:"minimum_reverification_interval,omitempty"`
}

// EditGeoSite updates settings of an existing Geo site.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#edit-a-geo-site
func (s *GeoSitesService) EditGeoSite(id int, opt *EditGeoSiteOptions, options ...RequestOptionFunc) (*GeoSite, *Response, error) {
	u := fmt.Sprintf("geo_sites/%d", id)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	site := new(GeoSite)
	resp, err := s.client.Do(req, site)
	if err != nil {
		return nil, resp, err
	}

	return site, resp, nil
}

// DeleteGeoSite removes the Geo site.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#delete-a-geo-site
func (s *GeoSitesService) DeleteGeoSite(id int, options ...RequestOptionFunc) (*Response, error) {
	u := fmt.Sprintf("geo_sites/%d", id)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RepairGeoSite to repair the OAuth authentication of a Geo site.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#repair-a-geo-site
func (s *GeoSitesService) RepairGeoSite(id int, options ...RequestOptionFunc) (*GeoSite, *Response, error) {
	u := fmt.Sprintf("geo_sites/%d/repair", id)

	req, err := s.client.NewRequest(http.MethodPost, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	site := new(GeoSite)
	resp, err := s.client.Do(req, site)
	if err != nil {
		return nil, resp, err
	}

	return site, resp, nil
}

// GeoSiteStatus represents the status of Geo Site.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#retrieve-status-about-all-geo-sites
type GeoSiteStatus struct {
	GeoNodeID                                       int               `json:"geo_node_id"`
	ProjectsCount                                   int               `json:"projects_count"`
	ContainerRepositoriesReplicationEnabled         bool              `json:"container_repositories_replication_enabled"`
	LFSObjectsCount                                 int               `json:"lfs_objects_count"`
	LFSObjectsChecksumTotalCount                    int               `json:"lfs_objects_checksum_total_count"`
	LFSObjectsChecksummedCount                      int               `json:"lfs_objects_checksummed_count"`
	LFSObjectsChecksumFailedCount                   int               `json:"lfs_objects_checksum_failed_count"`
	LFSObjectsSyncedCount                           int               `json:"lfs_objects_synced_count"`
	LFSObjectsFailedCount                           int               `json:"lfs_objects_failed_count"`
	LFSObjectsRegistryCount                         int               `json:"lfs_objects_registry_count"`
	LFSObjectsVerificationTotalCount                int               `json:"lfs_objects_verification_total_count"`
	LFSObjectsVerifiedCount                         int               `json:"lfs_objects_verified_count"`
	LFSObjectsVerificationFailedCount               int               `json:"lfs_objects_verification_failed_count"`
	MergeRequestDiffsCount                          int               `json:"merge_request_diffs_count"`
	MergeRequestDiffsChecksumTotalCount             int               `json:"merge_request_diffs_checksum_total_count"`
	MergeRequestDiffsChecksummedCount               int               `json:"merge_request_diffs_checksummed_count"`
	MergeRequestDiffsChecksumFailedCount            int               `json:"merge_request_diffs_checksum_failed_count"`
	MergeRequestDiffsSyncedCount                    int               `json:"merge_request_diffs_synced_count"`
	MergeRequestDiffsFailedCount                    int               `json:"merge_request_diffs_failed_count"`
	MergeRequestDiffsRegistryCount                  int               `json:"merge_request_diffs_registry_count"`
	MergeRequestDiffsVerificationTotalCount         int               `json:"merge_request_diffs_verification_total_count"`
	MergeRequestDiffsVerifiedCount                  int               `json:"merge_request_diffs_verified_count"`
	MergeRequestDiffsVerificationFailedCount        int               `json:"merge_request_diffs_verification_failed_count"`
	PackageFilesCount                               int               `json:"package_files_count"`
	PackageFilesChecksumTotalCount                  int               `json:"package_files_checksum_total_count"`
	PackageFilesChecksummedCount                    int               `json:"package_files_checksummed_count"`
	PackageFilesChecksumFailedCount                 int               `json:"package_files_checksum_failed_count"`
	PackageFilesSyncedCount                         int               `json:"package_files_synced_count"`
	PackageFilesFailedCount                         int               `json:"package_files_failed_count"`
	PackageFilesRegistryCount                       int               `json:"package_files_registry_count"`
	PackageFilesVerificationTotalCount              int               `json:"package_files_verification_total_count"`
	PackageFilesVerifiedCount                       int               `json:"package_files_verified_count"`
	PackageFilesVerificationFailedCount             int               `json:"package_files_verification_failed_count"`
	TerraformStateVersionsCount                     int               `json:"terraform_state_versions_count"`
	TerraformStateVersionsChecksumTotalCount        int               `json:"terraform_state_versions_checksum_total_count"`
	TerraformStateVersionsChecksummedCount          int               `json:"terraform_state_versions_checksummed_count"`
	TerraformStateVersionsChecksumFailedCount       int               `json:"terraform_state_versions_checksum_failed_count"`
	TerraformStateVersionsSyncedCount               int               `json:"terraform_state_versions_synced_count"`
	TerraformStateVersionsFailedCount               int               `json:"terraform_state_versions_failed_count"`
	TerraformStateVersionsRegistryCount             int               `json:"terraform_state_versions_registry_count"`
	TerraformStateVersionsVerificationTotalCount    int               `json:"terraform_state_versions_verification_total_count"`
	TerraformStateVersionsVerifiedCount             int               `json:"terraform_state_versions_verified_count"`
	TerraformStateVersionsVerificationFailedCount   int               `json:"terraform_state_versions_verification_failed_count"`
	SnippetRepositoriesCount                        int               `json:"snippet_repositories_count"`
	SnippetRepositoriesChecksumTotalCount           int               `json:"snippet_repositories_checksum_total_count"`
	SnippetRepositoriesChecksummedCount             int               `json:"snippet_repositories_checksummed_count"`
	SnippetRepositoriesChecksumFailedCount          int               `json:"snippet_repositories_checksum_failed_count"`
	SnippetRepositoriesSyncedCount                  int               `json:"snippet_repositories_synced_count"`
	SnippetRepositoriesFailedCount                  int               `json:"snippet_repositories_failed_count"`
	SnippetRepositoriesRegistryCount                int               `json:"snippet_repositories_registry_count"`
	SnippetRepositoriesVerificationTotalCount       int               `json:"snippet_repositories_verification_total_count"`
	SnippetRepositoriesVerifiedCount                int               `json:"snippet_repositories_verified_count"`
	SnippetRepositoriesVerificationFailedCount      int               `json:"snippet_repositories_verification_failed_count"`
	GroupWikiRepositoriesCount                      int               `json:"group_wiki_repositories_count"`
	GroupWikiRepositoriesChecksumTotalCount         int               `json:"group_wiki_repositories_checksum_total_count"`
	GroupWikiRepositoriesChecksummedCount           int               `json:"group_wiki_repositories_checksummed_count"`
	GroupWikiRepositoriesChecksumFailedCount        int               `json:"group_wiki_repositories_checksum_failed_count"`
	GroupWikiRepositoriesSyncedCount                int               `json:"group_wiki_repositories_synced_count"`
	GroupWikiRepositoriesFailedCount                int               `json:"group_wiki_repositories_failed_count"`
	GroupWikiRepositoriesRegistryCount              int               `json:"group_wiki_repositories_registry_count"`
	GrupWikiRepositoriesVerificationTotalCount      int               `json:"group_wiki_repositories_verification_total_count"`
	GroupWikiRepositoriesVerifiedCount              int               `json:"group_wiki_repositories_verified_count"`
	GroupWikiRepositoriesVerificationFailedCount    int               `json:"group_wiki_repositories_verification_failed_count"`
	PipelineArtifactsCount                          int               `json:"pipeline_artifacts_count"`
	PipelineArtifactsChecksumTotalCount             int               `json:"pipeline_artifacts_checksum_total_count"`
	PipelineArtifactsChecksummedCount               int               `json:"pipeline_artifacts_checksummed_count"`
	PipelineArtifactsChecksumFailedCount            int               `json:"pipeline_artifacts_checksum_failed_count"`
	PipelineArtifactsSyncedCount                    int               `json:"pipeline_artifacts_synced_count"`
	PipelineArtifactsFailedCount                    int               `json:"pipeline_artifacts_failed_count"`
	PipelineArtifactsRegistryCount                  int               `json:"pipeline_artifacts_registry_count"`
	PipelineArtifactsVerificationTotalCount         int               `json:"pipeline_artifacts_verification_total_count"`
	PipelineArtifactsVerifiedCount                  int               `json:"pipeline_artifacts_verified_count"`
	PipelineArtifactsVerificationFailedCount        int               `json:"pipeline_artifacts_verification_failed_count"`
	PagesDeploymentsCount                           int               `json:"pages_deployments_count"`
	PagesDeploymentsChecksumTotalCount              int               `json:"pages_deployments_checksum_total_count"`
	PagesDeploymentsChecksummedCount                int               `json:"pages_deployments_checksummed_count"`
	PagesDeploymentsChecksumFailedCount             int               `json:"pages_deployments_checksum_failed_count"`
	PagesDeploymentsSyncedCount                     int               `json:"pages_deployments_synced_count"`
	PagesDeploymentsFailedCount                     int               `json:"pages_deployments_failed_count"`
	PagesDeploymentsRegistryCount                   int               `json:"pages_deployments_registry_count"`
	PagesDeploymentsVerificationTotalCount          int               `json:"pages_deployments_verification_total_count"`
	PagesDeploymentsVerifiedCount                   int               `json:"pages_deployments_verified_count"`
	PagesDeploymentsVerificationFailedCount         int               `json:"pages_deployments_verification_failed_count"`
	UploadsCount                                    int               `json:"uploads_count"`
	UploadsChecksumTotalCount                       int               `json:"uploads_checksum_total_count"`
	UploadsChecksummedCount                         int               `json:"uploads_checksummed_count"`
	UploadsChecksumFailedCount                      int               `json:"uploads_checksum_failed_count"`
	UploadsSyncedCount                              int               `json:"uploads_synced_count"`
	UploadsFailedCount                              int               `json:"uploads_failed_count"`
	UploadsRegistryCount                            int               `json:"uploads_registry_count"`
	UploadsVerificationTotalCount                   int               `json:"uploads_verification_total_count"`
	UploadsVerifiedCount                            int               `json:"uploads_verified_count"`
	UploadsVerificationFailedCount                  int               `json:"uploads_verification_failed_count"`
	JobArtifactsCount                               int               `json:"job_artifacts_count"`
	JobArtifactsChecksumTotalCount                  int               `json:"job_artifacts_checksum_total_count"`
	JobArtifactsChecksummedCount                    int               `json:"job_artifacts_checksummed_count"`
	JobArtifactsChecksumFailedCount                 int               `json:"job_artifacts_checksum_failed_count"`
	JobArtifactsSyncedCount                         int               `json:"job_artifacts_synced_count"`
	JobArtifactsFailedCount                         int               `json:"job_artifacts_failed_count"`
	JobArtifactsRegistryCount                       int               `json:"job_artifacts_registry_count"`
	JobArtifactsVerificationTotalCount              int               `json:"job_artifacts_verification_total_count"`
	JobArtifactsVerifiedCount                       int               `json:"job_artifacts_verified_count"`
	JobArtifactsVerificationFailedCount             int               `json:"job_artifacts_verification_failed_count"`
	CISecureFilesCount                              int               `json:"ci_secure_files_count"`
	CISecureFilesChecksumTotalCount                 int               `json:"ci_secure_files_checksum_total_count"`
	CISecureFilesChecksummedCount                   int               `json:"ci_secure_files_checksummed_count"`
	CISecureFilesChecksumFailedCount                int               `json:"ci_secure_files_checksum_failed_count"`
	CISecureFilesSyncedCount                        int               `json:"ci_secure_files_synced_count"`
	CISecureFilesFailedCount                        int               `json:"ci_secure_files_failed_count"`
	CISecureFilesRegistryCount                      int               `json:"ci_secure_files_registry_count"`
	CISecureFilesVerificationTotalCount             int               `json:"ci_secure_files_verification_total_count"`
	CISecureFilesVerifiedCount                      int               `json:"ci_secure_files_verified_count"`
	CISecureFilesVerificationFailedCount            int               `json:"ci_secure_files_verification_failed_count"`
	ContainerRepositoriesCount                      int               `json:"container_repositories_count"`
	ContainerRepositoriesChecksumTotalCount         int               `json:"container_repositories_checksum_total_count"`
	ContainerRepositoriesChecksummedCount           int               `json:"container_repositories_checksummed_count"`
	ContainerRepositoriesChecksumFailedCount        int               `json:"container_repositories_checksum_failed_count"`
	ContainerRepositoriesSyncedCount                int               `json:"container_repositories_synced_count"`
	ContainerRepositoriesFailedCount                int               `json:"container_repositories_failed_count"`
	ContainerRepositoriesRegistryCount              int               `json:"container_repositories_registry_count"`
	ContainerRepositoriesVerificationTotalCount     int               `json:"container_repositories_verification_total_count"`
	ContainerRepositoriesVerifiedCount              int               `json:"container_repositories_verified_count"`
	ContainerRepositoriesVerificationFailedCount    int               `json:"container_repositories_verification_failed_count"`
	DependencyProxyBlobsCount                       int               `json:"dependency_proxy_blobs_count"`
	DependencyProxyBlobsChecksumTotalCount          int               `json:"dependency_proxy_blobs_checksum_total_count"`
	DependencyProxyBlobsChecksummedCount            int               `json:"dependency_proxy_blobs_checksummed_count"`
	DependencyProxyBlobsChecksumFailedCount         int               `json:"dependency_proxy_blobs_checksum_failed_count"`
	DependencyProxyBlobsSyncedCount                 int               `json:"dependency_proxy_blobs_synced_count"`
	DependencyProxyBlobsFailedCount                 int               `json:"dependency_proxy_blobs_failed_count"`
	DependencyProxyBlobsRegistryCount               int               `json:"dependency_proxy_blobs_registry_count"`
	DependencyProxyBlobsVerificationTotalCount      int               `json:"dependency_proxy_blobs_verification_total_count"`
	DependencyProxyBlobsVerifiedCount               int               `json:"dependency_proxy_blobs_verified_count"`
	DependencyProxyBlobsVerificationFailedCount     int               `json:"dependency_proxy_blobs_verification_failed_count"`
	DependencyProxyManifestsCount                   int               `json:"dependency_proxy_manifests_count"`
	DependencyProxyManifestsChecksumTotalCount      int               `json:"dependency_proxy_manifests_checksum_total_count"`
	DependencyProxyManifestsChecksummedCount        int               `json:"dependency_proxy_manifests_checksummed_count"`
	DependencyProxyManifestsChecksumFailedCount     int               `json:"dependency_proxy_manifests_checksum_failed_count"`
	DependencyProxyManifestsSyncedCount             int               `json:"dependency_proxy_manifests_synced_count"`
	DependencyProxyManifestsFailedCount             int               `json:"dependency_proxy_manifests_failed_count"`
	DependencyProxyManifestsRegistryCount           int               `json:"dependency_proxy_manifests_registry_count"`
	DependencyProxyManifestsVerificationTotalCount  int               `json:"dependency_proxy_manifests_verification_total_count"`
	DependencyProxyManifestsVerifiedCount           int               `json:"dependency_proxy_manifests_verified_count"`
	DependencyProxyManifestsVerificationFailedCount int               `json:"dependency_proxy_manifests_verification_failed_count"`
	ProjectWikiRepositoriesCount                    int               `json:"project_wiki_repositories_count"`
	ProjectWikiRepositoriesChecksumTotalCount       int               `json:"project_wiki_repositories_checksum_total_count"`
	ProjectWikiRepositoriesChecksummedCount         int               `json:"project_wiki_repositories_checksummed_count"`
	ProjectWikiRepositoriesChecksumFailedCount      int               `json:"project_wiki_repositories_checksum_failed_count"`
	ProjectWikiRepositoriesSyncedCount              int               `json:"project_wiki_repositories_synced_count"`
	ProjectWikiRepositoriesFailedCount              int               `json:"project_wiki_repositories_failed_count"`
	ProjectWikiRepositoriesRegistryCount            int               `json:"project_wiki_repositories_registry_count"`
	ProjectWikiRepositoriesVerificationTotalCount   int               `json:"project_wiki_repositories_verification_total_count"`
	ProjectWikiRepositoriesVerifiedCount            int               `json:"project_wiki_repositories_verified_count"`
	ProjectWikiRepositoriesVerificationFailedCount  int               `json:"project_wiki_repositories_verification_failed_count"`
	GitFetchEventCountWeekly                        int               `json:"git_fetch_event_count_weekly"`
	GitPushEventCountWeekly                         int               `json:"git_push_event_count_weekly"`
	ProxyRemoteRequestsEventCountWeekly             int               `json:"proxy_remote_requests_event_count_weekly"`
	ProxyLocalRequestsEventCountWeekly              int               `json:"proxy_local_requests_event_count_weekly"`
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
	ReplicationSlotsCount                           int               `json:"replication_slots_count"`
	ReplicationSlotsUsedCount                       int               `json:"replication_slots_used_count"`
	Healthy                                         bool              `json:"healthy"`
	Health                                          string            `json:"health"`
	HealthStatus                                    string            `json:"health_status"`
	MissingOAuthApplication                         bool              `json:"missing_oauth_application"`
	DBReplicationLagSeconds                         int               `json:"db_replication_lag_seconds"`
	ReplicationSlotsMaxRetainedWalBytes             int               `json:"replication_slots_max_retained_wal_bytes"`
	RepositoriesCheckedCount                        int               `json:"repositories_checked_count"`
	RepositoriesCheckedFailedCount                  int               `json:"repositories_checked_failed_count"`
	LastEventID                                     int               `json:"last_event_id"`
	LastEventTimestamp                              int               `json:"last_event_timestamp"`
	CursorLastEventID                               int               `json:"cursor_last_event_id"`
	CursorLastEventTimestamp                        int               `json:"cursor_last_event_timestamp"`
	LastSuccessfulStatusCheckTimestamp              int               `json:"last_successful_status_check_timestamp"`
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
// https://docs.gitlab.com/api/geo_sites/#retrieve-status-about-all-geo-sites
type GeoSiteStatusLink struct {
	Self string `json:"self"`
	Site string `json:"site"`
}

// ListStatusOfAllGeoSitesOptions represents the available ListStatusOfAllGeoSites() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#retrieve-status-about-all-geo-sites
type ListStatusOfAllGeoSitesOptions ListOptions

// ListStatusOfAllGeoSites get the list of status of all Geo Sites.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#retrieve-status-about-all-geo-sites
func (s *GeoSitesService) ListStatusOfAllGeoSites(opt *ListStatusOfAllGeoSitesOptions, options ...RequestOptionFunc) ([]*GeoSiteStatus, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "geo_sites/status", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var statuses []*GeoSiteStatus
	resp, err := s.client.Do(req, &statuses)
	if err != nil {
		return nil, resp, err
	}

	return statuses, resp, nil
}

// GetStatusOfGeoSite get the of status of a specific Geo Sites.
//
// GitLab API docs:
// https://docs.gitlab.com/api/geo_sites/#retrieve-status-about-a-specific-geo-site
func (s *GeoSitesService) GetStatusOfGeoSite(id int, options ...RequestOptionFunc) (*GeoSiteStatus, *Response, error) {
	u := fmt.Sprintf("geo_sites/%d/status", id)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	status := new(GeoSiteStatus)
	resp, err := s.client.Do(req, status)
	if err != nil {
		return nil, resp, err
	}

	return status, resp, nil
}
