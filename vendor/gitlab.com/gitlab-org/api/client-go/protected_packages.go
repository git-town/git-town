package gitlab

import (
	"fmt"
	"net/http"
)

type (
	ProtectedPackagesServiceInterface interface {
		// ListPackageProtectionRules gets a list of project package protection rules.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_packages_protection_rules/#list-package-protection-rules
		ListPackageProtectionRules(pid any, opt *ListPackageProtectionRulesOptions, options ...RequestOptionFunc) ([]*PackageProtectionRule, *Response, error)
		// CreatePackageProtectionRules creates a new package protection rules.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_packages_protection_rules/#create-a-package-protection-rule
		CreatePackageProtectionRules(pid any, opt *CreatePackageProtectionRulesOptions, options ...RequestOptionFunc) (*PackageProtectionRule, *Response, error)
		// UpdatePackageProtectionRules updates an existing package protection rule.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_packages_protection_rules/#update-a-package-protection-rule
		UpdatePackageProtectionRules(pid any, packageProtectionRule int64, opt *UpdatePackageProtectionRulesOptions, options ...RequestOptionFunc) (*PackageProtectionRule, *Response, error)
		// DeletePackageProtectionRules deletes an existing package protection rules.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_packages_protection_rules/#delete-a-package-protection-rule
		DeletePackageProtectionRules(pid any, packageProtectionRule int64, options ...RequestOptionFunc) (*Response, error)
	}

	// ProtectedPackagesService handles communication with the protected packages related methods
	// of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/project_packages_protection_rules/
	ProtectedPackagesService struct {
		client *Client
	}
)

var _ ProtectedPackagesServiceInterface = (*ProtectedPackagesService)(nil)

// PackageProtectionRule represents a GitLab package protection rule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_packages_protection_rules
type PackageProtectionRule struct {
	ID                          int64  `json:"id"`
	ProjectID                   int64  `json:"project_id"`
	PackageNamePattern          string `json:"package_name_pattern"`
	PackageType                 string `json:"package_type"`
	MinimumAccessLevelForDelete string `json:"minimum_access_level_for_delete"`
	MinimumAccessLevelForPush   string `json:"minimum_access_level_for_push"`
}

// ListPackageProtectionRulesOptions represents the available ListPackageProtectionRules() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_packages_protection_rules/#list-package-protection-rules
type ListPackageProtectionRulesOptions struct {
	ListOptions
}

// CreatePackageProtectionRulesOptions represents the available CreatePackageProtectionRules() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_packages_protection_rules/#create-a-package-protection-rule
type CreatePackageProtectionRulesOptions struct {
	PackageNamePattern          *string `url:"package_name_pattern" json:"package_name_pattern"`
	PackageType                 *string `url:"package_type" json:"package_type"`
	MinimumAccessLevelForDelete *string `url:"minimum_access_level_for_delete" json:"minimum_access_level_for_delete"`
	MinimumAccessLevelForPush   *string `url:"minimum_access_level_for_push" json:"minimum_access_level_for_push"`
}

// UpdatePackageProtectionRulesOptions represents the available
// UpdatePackageProtectionRules() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_packages_protection_rules/#update-a-package-protection-rule
type UpdatePackageProtectionRulesOptions struct {
	PackageNamePattern          *string `url:"package_name_pattern" json:"package_name_pattern"`
	PackageType                 *string `url:"package_type" json:"package_type"`
	MinimumAccessLevelForDelete *string `url:"minimum_access_level_for_delete" json:"minimum_access_level_for_delete"`
	MinimumAccessLevelForPush   *string `url:"minimum_access_level_for_push" json:"minimum_access_level_for_push"`
}

func (s *ProtectedPackagesService) ListPackageProtectionRules(pid any, opts *ListPackageProtectionRulesOptions, options ...RequestOptionFunc) ([]*PackageProtectionRule, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}

	u := fmt.Sprintf("projects/%s/packages/protection/rules", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opts, options)
	if err != nil {
		return nil, nil, err
	}

	var b []*PackageProtectionRule
	resp, err := s.client.Do(req, &b)
	if err != nil {
		return nil, resp, err
	}

	return b, resp, nil
}

func (s *ProtectedPackagesService) CreatePackageProtectionRules(pid any, opt *CreatePackageProtectionRulesOptions, options ...RequestOptionFunc) (*PackageProtectionRule, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/packages/protection/rules", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	p := new(PackageProtectionRule)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *ProtectedPackagesService) DeletePackageProtectionRules(pid any, packageProtectionRule int64, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/packages/protection/rules/%d", PathEscape(project), packageProtectionRule)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

func (s *ProtectedPackagesService) UpdatePackageProtectionRules(pid any, packageProtectionRule int64, opt *UpdatePackageProtectionRulesOptions, options ...RequestOptionFunc) (*PackageProtectionRule, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/packages/protection/rules/%d", PathEscape(project), packageProtectionRule)

	req, err := s.client.NewRequest(http.MethodPatch, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	p := new(PackageProtectionRule)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}
