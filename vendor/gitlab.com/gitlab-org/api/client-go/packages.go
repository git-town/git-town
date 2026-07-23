//
// Copyright 2021, Kordian Bruck
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
	PackagesServiceInterface interface {
		ListProjectPackages(pid any, opt *ListProjectPackagesOptions, options ...RequestOptionFunc) ([]*Package, *Response, error)
		ListGroupPackages(gid any, opt *ListGroupPackagesOptions, options ...RequestOptionFunc) ([]*GroupPackage, *Response, error)
		ListPackageFiles(pid any, pkg int64, opt *ListPackageFilesOptions, options ...RequestOptionFunc) ([]*PackageFile, *Response, error)
		DeleteProjectPackage(pid any, pkg int64, options ...RequestOptionFunc) (*Response, error)
		DeletePackageFile(pid any, pkg, file int64, options ...RequestOptionFunc) (*Response, error)
	}

	// PackagesService handles communication with the packages related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/packages/
	PackagesService struct {
		client *Client
	}
)

var _ PackagesServiceInterface = (*PackagesService)(nil)

// Package represents a GitLab package.
//
// GitLab API docs: https://docs.gitlab.com/api/packages/
type Package struct {
	ID               int64         `json:"id"`
	Name             string        `json:"name"`
	Version          string        `json:"version"`
	PackageType      string        `json:"package_type"`
	Status           string        `json:"status"`
	Links            *PackageLinks `json:"_links"`
	CreatedAt        *time.Time    `json:"created_at"`
	LastDownloadedAt *time.Time    `json:"last_downloaded_at"`
	Tags             []PackageTag  `json:"tags"`
}

func (s Package) String() string {
	return Stringify(s)
}

// GroupPackage represents a GitLab group package.
//
// GitLab API docs: https://docs.gitlab.com/api/packages/
type GroupPackage struct {
	Package
	ProjectID   int64  `json:"project_id"`
	ProjectPath string `json:"project_path"`
}

func (s GroupPackage) String() string {
	return Stringify(s)
}

// PackageLinks holds links for itself and deleting.
type PackageLinks struct {
	WebPath       string `json:"web_path"`
	DeleteAPIPath string `json:"delete_api_path"`
}

func (s PackageLinks) String() string {
	return Stringify(s)
}

// PackageTag holds label information about the package
type PackageTag struct {
	ID        int64      `json:"id"`
	PackageID int64      `json:"package_id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (s PackageTag) String() string {
	return Stringify(s)
}

// PackageFile represents one file contained within a package.
//
// GitLab API docs: https://docs.gitlab.com/api/packages/
type PackageFile struct {
	ID         int64       `json:"id"`
	PackageID  int64       `json:"package_id"`
	CreatedAt  *time.Time  `json:"created_at"`
	FileName   string      `json:"file_name"`
	Size       int64       `json:"size"`
	FileMD5    string      `json:"file_md5"`
	FileSHA1   string      `json:"file_sha1"`
	FileSHA256 string      `json:"file_sha256"`
	Pipeline   *[]Pipeline `json:"pipelines"`
}

func (s PackageFile) String() string {
	return Stringify(s)
}

// ListProjectPackagesOptions represents the available ListProjectPackages()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/packages/#for-a-project
type ListProjectPackagesOptions struct {
	ListOptions
	OrderBy            *string `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort               *string `url:"sort,omitempty" json:"sort,omitempty"`
	PackageType        *string `url:"package_type,omitempty" json:"package_type,omitempty"`
	PackageName        *string `url:"package_name,omitempty" json:"package_name,omitempty"`
	PackageVersion     *string `url:"package_version,omitempty" json:"package_version,omitempty"`
	IncludeVersionless *bool   `url:"include_versionless,omitempty" json:"include_versionless,omitempty"`
	Status             *string `url:"status,omitempty" json:"status,omitempty"`
}

// ListProjectPackages gets a list of packages in a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/packages/#for-a-project
func (s *PackagesService) ListProjectPackages(pid any, opt *ListProjectPackagesOptions, options ...RequestOptionFunc) ([]*Package, *Response, error) {
	return do[[]*Package](s.client,
		withPath("projects/%s/packages", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListGroupPackagesOptions represents the available ListGroupPackages()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/packages/#for-a-group
type ListGroupPackagesOptions struct {
	ListOptions
	ExcludeSubGroups   *bool   `url:"exclude_subgroups,omitempty" json:"exclude_subgroups,omitempty"`
	OrderBy            *string `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort               *string `url:"sort,omitempty" json:"sort,omitempty"`
	PackageType        *string `url:"package_type,omitempty" json:"package_type,omitempty"`
	PackageName        *string `url:"package_name,omitempty" json:"package_name,omitempty"`
	IncludeVersionless *bool   `url:"include_versionless,omitempty" json:"include_versionless,omitempty"`
	Status             *string `url:"status,omitempty" json:"status,omitempty"`
}

// ListGroupPackages gets a list of packages in a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/packages/#for-a-group
func (s *PackagesService) ListGroupPackages(gid any, opt *ListGroupPackagesOptions, options ...RequestOptionFunc) ([]*GroupPackage, *Response, error) {
	return do[[]*GroupPackage](s.client,
		withPath("groups/%s/packages", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListPackageFilesOptions represents the available ListPackageFiles()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/packages/#list-package-files
type ListPackageFilesOptions struct {
	ListOptions
}

// ListPackageFiles gets a list of files that are within a package
//
// GitLab API docs:
// https://docs.gitlab.com/api/packages/#list-package-files
func (s *PackagesService) ListPackageFiles(pid any, pkg int64, opt *ListPackageFilesOptions, options ...RequestOptionFunc) ([]*PackageFile, *Response, error) {
	return do[[]*PackageFile](s.client,
		withPath("projects/%s/packages/%d/package_files", ProjectID{pid}, pkg),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteProjectPackage deletes a package in a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/packages/#delete-a-project-package
func (s *PackagesService) DeleteProjectPackage(pid any, pkg int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/packages/%d", ProjectID{pid}, pkg),
		withRequestOpts(options...),
	)
	return resp, err
}

// DeletePackageFile deletes a file in project package
//
// GitLab API docs:
// https://docs.gitlab.com/api/packages/#delete-a-package-file
func (s *PackagesService) DeletePackageFile(pid any, pkg, file int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/packages/%d/package_files/%d", ProjectID{pid}, pkg, file),
		withRequestOpts(options...),
	)
	return resp, err
}
