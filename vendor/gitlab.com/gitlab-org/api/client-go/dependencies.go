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

type (
	DependenciesServiceInterface interface {
		// ListProjectDependencies Get a list of project dependencies. This API partially
		// mirroring Dependency List feature. This list can be generated only for languages
		// and package managers supported by Gemnasium.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/dependencies/#list-project-dependencies
		ListProjectDependencies(pid any, opt *ListProjectDependenciesOptions, options ...RequestOptionFunc) ([]*Dependency, *Response, error)
	}

	// DependenciesService handles communication with the dependencies related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/dependencies/
	DependenciesService struct {
		client *Client
	}
)

var _ DependenciesServiceInterface = (*DependenciesService)(nil)

// Dependency represents a project dependency.
//
// GitLab API docs: https://docs.gitlab.com/api/dependencies/
type Dependency struct {
	Name               string                        `url:"name" json:"name"`
	Version            string                        `url:"version" json:"version"`
	PackageManager     DependencyPackageManagerValue `url:"package_manager" json:"package_manager"`
	DependencyFilePath string                        `url:"dependency_file_path" json:"dependency_file_path"`
	Vulnerabilities    []*DependencyVulnerability    `url:"vulnerabilities" json:"vulnerabilities"`
	Licenses           []*DependencyLicense          `url:"licenses" json:"licenses"`
}

// DependencyVulnerability represents a project dependency vulnerability.
//
// GitLab API docs: https://docs.gitlab.com/api/dependencies/
type DependencyVulnerability struct {
	Name     string `url:"name" json:"name"`
	Severity string `url:"severity" json:"severity"`
	ID       int64  `url:"id" json:"id"`
	URL      string `url:"url" json:"url"`
}

// DependencyLicense represents a project dependency license.
//
// GitLab API docs: https://docs.gitlab.com/api/dependencies/
type DependencyLicense struct {
	Name string `url:"name" json:"name"`
	URL  string `url:"url" json:"url"`
}

// ListProjectDependenciesOptions represents the options for listing project
// dependencies.
//
// GitLab API docs:
// https://docs.gitlab.com/api/dependencies/#list-project-dependencies
type ListProjectDependenciesOptions struct {
	ListOptions
	PackageManager []*DependencyPackageManagerValue `url:"package_manager,comma,omitempty" json:"package_manager,omitempty"`
}

func (s *DependenciesService) ListProjectDependencies(pid any, opt *ListProjectDependenciesOptions, options ...RequestOptionFunc) ([]*Dependency, *Response, error) {
	return do[[]*Dependency](s.client,
		withPath("projects/%s/dependencies", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
