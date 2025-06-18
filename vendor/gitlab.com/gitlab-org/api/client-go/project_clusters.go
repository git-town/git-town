//
// Copyright 2021, Matej Velikonja
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
	// Deprecated: in GitLab 14.5, to be removed in 19.0
	ProjectClustersServiceInterface interface {
		// Deprecated: in GitLab 14.5, to be removed in 19.0
		ListClusters(pid any, options ...RequestOptionFunc) ([]*ProjectCluster, *Response, error)
		// Deprecated: in GitLab 14.5, to be removed in 19.0
		GetCluster(pid any, cluster int, options ...RequestOptionFunc) (*ProjectCluster, *Response, error)
		// Deprecated: in GitLab 14.5, to be removed in 19.0
		AddCluster(pid any, opt *AddClusterOptions, options ...RequestOptionFunc) (*ProjectCluster, *Response, error)
		// Deprecated: in GitLab 14.5, to be removed in 19.0
		EditCluster(pid any, cluster int, opt *EditClusterOptions, options ...RequestOptionFunc) (*ProjectCluster, *Response, error)
		// Deprecated: in GitLab 14.5, to be removed in 19.0
		DeleteCluster(pid any, cluster int, options ...RequestOptionFunc) (*Response, error)
	}

	// ProjectClustersService handles communication with the
	// project clusters related methods of the GitLab API.
	// Deprecated: in GitLab 14.5, to be removed in 19.0
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/project_clusters/
	ProjectClustersService struct {
		client *Client
	}
)

// Deprecated: in GitLab 14.5, to be removed in 19.0
var _ ProjectClustersServiceInterface = (*ProjectClustersService)(nil)

// ProjectCluster represents a GitLab Project Cluster.
// Deprecated: in GitLab 14.5, to be removed in 19.0
//
// GitLab API docs: https://docs.gitlab.com/api/project_clusters/
type ProjectCluster struct {
	ID                 int                 `json:"id"`
	Name               string              `json:"name"`
	Domain             string              `json:"domain"`
	CreatedAt          *time.Time          `json:"created_at"`
	ProviderType       string              `json:"provider_type"`
	PlatformType       string              `json:"platform_type"`
	EnvironmentScope   string              `json:"environment_scope"`
	ClusterType        string              `json:"cluster_type"`
	User               *User               `json:"user"`
	PlatformKubernetes *PlatformKubernetes `json:"platform_kubernetes"`
	ManagementProject  *ManagementProject  `json:"management_project"`
	Project            *Project            `json:"project"`
}

// Deprecated: in GitLab 14.5, to be removed in 19.0
func (v ProjectCluster) String() string {
	return Stringify(v)
}

// PlatformKubernetes represents a GitLab Project Cluster PlatformKubernetes.
// Deprecated: in GitLab 14.5, to be removed in 19.0
type PlatformKubernetes struct {
	APIURL            string `json:"api_url"`
	Token             string `json:"token"`
	CaCert            string `json:"ca_cert"`
	Namespace         string `json:"namespace"`
	AuthorizationType string `json:"authorization_type"`
}

// ManagementProject represents a GitLab Project Cluster management_project.
// Deprecated: in GitLab 14.5, to be removed in 19.0
type ManagementProject struct {
	ID                int        `json:"id"`
	Description       string     `json:"description"`
	Name              string     `json:"name"`
	NameWithNamespace string     `json:"name_with_namespace"`
	Path              string     `json:"path"`
	PathWithNamespace string     `json:"path_with_namespace"`
	CreatedAt         *time.Time `json:"created_at"`
}

// ListClusters gets a list of all clusters in a project.
// Deprecated: in GitLab 14.5, to be removed in 19.0
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_clusters/#list-project-clusters
func (s *ProjectClustersService) ListClusters(pid any, options ...RequestOptionFunc) ([]*ProjectCluster, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/clusters", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var pcs []*ProjectCluster
	resp, err := s.client.Do(req, &pcs)
	if err != nil {
		return nil, resp, err
	}

	return pcs, resp, nil
}

// GetCluster gets a cluster.
// Deprecated: in GitLab 14.5, to be removed in 19.0
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_clusters/#get-a-single-project-cluster
func (s *ProjectClustersService) GetCluster(pid any, cluster int, options ...RequestOptionFunc) (*ProjectCluster, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/clusters/%d", PathEscape(project), cluster)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	pc := new(ProjectCluster)
	resp, err := s.client.Do(req, &pc)
	if err != nil {
		return nil, resp, err
	}

	return pc, resp, nil
}

// AddClusterOptions represents the available AddCluster() options.
// Deprecated: in GitLab 14.5, to be removed in 19.0
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_clusters/#add-existing-cluster-to-project
type AddClusterOptions struct {
	Name                *string                       `url:"name,omitempty" json:"name,omitempty"`
	Domain              *string                       `url:"domain,omitempty" json:"domain,omitempty"`
	Enabled             *bool                         `url:"enabled,omitempty" json:"enabled,omitempty"`
	Managed             *bool                         `url:"managed,omitempty" json:"managed,omitempty"`
	EnvironmentScope    *string                       `url:"environment_scope,omitempty" json:"environment_scope,omitempty"`
	PlatformKubernetes  *AddPlatformKubernetesOptions `url:"platform_kubernetes_attributes,omitempty" json:"platform_kubernetes_attributes,omitempty"`
	ManagementProjectID *string                       `url:"management_project_id,omitempty" json:"management_project_id,omitempty"`
}

// AddPlatformKubernetesOptions represents the available PlatformKubernetes options for adding.
// Deprecated: in GitLab 14.5, to be removed in 19.0
type AddPlatformKubernetesOptions struct {
	APIURL            *string `url:"api_url,omitempty" json:"api_url,omitempty"`
	Token             *string `url:"token,omitempty" json:"token,omitempty"`
	CaCert            *string `url:"ca_cert,omitempty" json:"ca_cert,omitempty"`
	Namespace         *string `url:"namespace,omitempty" json:"namespace,omitempty"`
	AuthorizationType *string `url:"authorization_type,omitempty" json:"authorization_type,omitempty"`
}

// AddCluster adds an existing cluster to the project.
// Deprecated: in GitLab 14.5, to be removed in 19.0
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_clusters/#add-existing-cluster-to-project
func (s *ProjectClustersService) AddCluster(pid any, opt *AddClusterOptions, options ...RequestOptionFunc) (*ProjectCluster, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/clusters/user", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	pc := new(ProjectCluster)
	resp, err := s.client.Do(req, pc)
	if err != nil {
		return nil, resp, err
	}

	return pc, resp, nil
}

// EditClusterOptions represents the available EditCluster() options.
// Deprecated: in GitLab 14.5, to be removed in 19.0
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_clusters/#edit-project-cluster
type EditClusterOptions struct {
	Name                *string                        `url:"name,omitempty" json:"name,omitempty"`
	Domain              *string                        `url:"domain,omitempty" json:"domain,omitempty"`
	EnvironmentScope    *string                        `url:"environment_scope,omitempty" json:"environment_scope,omitempty"`
	ManagementProjectID *string                        `url:"management_project_id,omitempty" json:"management_project_id,omitempty"`
	PlatformKubernetes  *EditPlatformKubernetesOptions `url:"platform_kubernetes_attributes,omitempty" json:"platform_kubernetes_attributes,omitempty"`
}

// EditPlatformKubernetesOptions represents the available PlatformKubernetes options for editing.
// Deprecated: in GitLab 14.5, to be removed in 19.0
type EditPlatformKubernetesOptions struct {
	APIURL    *string `url:"api_url,omitempty" json:"api_url,omitempty"`
	Token     *string `url:"token,omitempty" json:"token,omitempty"`
	CaCert    *string `url:"ca_cert,omitempty" json:"ca_cert,omitempty"`
	Namespace *string `url:"namespace,omitempty" json:"namespace,omitempty"`
}

// EditCluster updates an existing project cluster.
// Deprecated: in GitLab 14.5, to be removed in 19.0
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_clusters/#edit-project-cluster
func (s *ProjectClustersService) EditCluster(pid any, cluster int, opt *EditClusterOptions, options ...RequestOptionFunc) (*ProjectCluster, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/clusters/%d", PathEscape(project), cluster)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	pc := new(ProjectCluster)
	resp, err := s.client.Do(req, pc)
	if err != nil {
		return nil, resp, err
	}

	return pc, resp, nil
}

// DeleteCluster deletes an existing project cluster.
// Deprecated: in GitLab 14.5, to be removed in 19.0
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_clusters/#delete-project-cluster
func (s *ProjectClustersService) DeleteCluster(pid any, cluster int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/clusters/%d", PathEscape(project), cluster)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
