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

import (
	"net/http"
)

type (
	CustomAttributesServiceInterface interface {
		// ListCustomUserAttributes lists the custom attributes of the specified user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#list-custom-attributes
		ListCustomUserAttributes(user int64, options ...RequestOptionFunc) ([]*CustomAttribute, *Response, error)

		// ListCustomGroupAttributes lists the custom attributes of the specified group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#list-custom-attributes
		ListCustomGroupAttributes(group int64, options ...RequestOptionFunc) ([]*CustomAttribute, *Response, error)

		// ListCustomProjectAttributes lists the custom attributes of the specified project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#list-custom-attributes
		ListCustomProjectAttributes(project int64, options ...RequestOptionFunc) ([]*CustomAttribute, *Response, error)

		// GetCustomUserAttribute returns the user attribute with a specific key.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#single-custom-attribute
		GetCustomUserAttribute(user int64, key string, options ...RequestOptionFunc) (*CustomAttribute, *Response, error)

		// GetCustomGroupAttribute returns the group attribute with a specific key.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#single-custom-attribute
		GetCustomGroupAttribute(group int64, key string, options ...RequestOptionFunc) (*CustomAttribute, *Response, error)

		// GetCustomProjectAttribute returns the project attribute with a specific key.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#single-custom-attribute
		GetCustomProjectAttribute(project int64, key string, options ...RequestOptionFunc) (*CustomAttribute, *Response, error)

		// SetCustomUserAttribute sets the custom attributes of the specified user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#set-custom-attribute
		SetCustomUserAttribute(user int64, c CustomAttribute, options ...RequestOptionFunc) (*CustomAttribute, *Response, error)

		// SetCustomGroupAttribute sets the custom attributes of the specified group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#set-custom-attribute
		SetCustomGroupAttribute(group int64, c CustomAttribute, options ...RequestOptionFunc) (*CustomAttribute, *Response, error)

		// SetCustomProjectAttribute sets the custom attributes of the specified project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#set-custom-attribute
		SetCustomProjectAttribute(project int64, c CustomAttribute, options ...RequestOptionFunc) (*CustomAttribute, *Response, error)

		// DeleteCustomUserAttribute removes the custom attribute of the specified user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#delete-custom-attribute
		DeleteCustomUserAttribute(user int64, key string, options ...RequestOptionFunc) (*Response, error)

		// DeleteCustomGroupAttribute removes the custom attribute of the specified group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#delete-custom-attribute
		DeleteCustomGroupAttribute(group int64, key string, options ...RequestOptionFunc) (*Response, error)

		// DeleteCustomProjectAttribute removes the custom attribute of the specified project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/custom_attributes/#delete-custom-attribute
		DeleteCustomProjectAttribute(project int64, key string, options ...RequestOptionFunc) (*Response, error)
	}

	// CustomAttributesService handles communication with the group, project and
	// user custom attributes related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/custom_attributes/
	CustomAttributesService struct {
		client *Client
	}
)

var _ CustomAttributesServiceInterface = (*CustomAttributesService)(nil)

// CustomAttribute struct is used to unmarshal response to api calls.
//
// GitLab API docs: https://docs.gitlab.com/api/custom_attributes/
type CustomAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *CustomAttributesService) ListCustomUserAttributes(user int64, options ...RequestOptionFunc) ([]*CustomAttribute, *Response, error) {
	return s.listCustomAttributes("users", user, options...)
}

func (s *CustomAttributesService) ListCustomGroupAttributes(group int64, options ...RequestOptionFunc) ([]*CustomAttribute, *Response, error) {
	return s.listCustomAttributes("groups", group, options...)
}

func (s *CustomAttributesService) ListCustomProjectAttributes(project int64, options ...RequestOptionFunc) ([]*CustomAttribute, *Response, error) {
	return s.listCustomAttributes("projects", project, options...)
}

func (s *CustomAttributesService) listCustomAttributes(resource string, id int64, options ...RequestOptionFunc) ([]*CustomAttribute, *Response, error) {
	res, resp, err := do[[]*CustomAttribute](s.client,
		withMethod(http.MethodGet),
		withPath("%s/%d/custom_attributes", resource, id),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

func (s *CustomAttributesService) GetCustomUserAttribute(user int64, key string, options ...RequestOptionFunc) (*CustomAttribute, *Response, error) {
	return s.getCustomAttribute("users", user, key, options...)
}

func (s *CustomAttributesService) GetCustomGroupAttribute(group int64, key string, options ...RequestOptionFunc) (*CustomAttribute, *Response, error) {
	return s.getCustomAttribute("groups", group, key, options...)
}

func (s *CustomAttributesService) GetCustomProjectAttribute(project int64, key string, options ...RequestOptionFunc) (*CustomAttribute, *Response, error) {
	return s.getCustomAttribute("projects", project, key, options...)
}

func (s *CustomAttributesService) getCustomAttribute(resource string, id int64, key string, options ...RequestOptionFunc) (*CustomAttribute, *Response, error) {
	res, resp, err := do[*CustomAttribute](s.client,
		withMethod(http.MethodGet),
		withPath("%s/%d/custom_attributes/%s", resource, id, key),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

func (s *CustomAttributesService) SetCustomUserAttribute(user int64, c CustomAttribute, options ...RequestOptionFunc) (*CustomAttribute, *Response, error) {
	return s.setCustomAttribute("users", user, c, options...)
}

func (s *CustomAttributesService) SetCustomGroupAttribute(group int64, c CustomAttribute, options ...RequestOptionFunc) (*CustomAttribute, *Response, error) {
	return s.setCustomAttribute("groups", group, c, options...)
}

func (s *CustomAttributesService) SetCustomProjectAttribute(project int64, c CustomAttribute, options ...RequestOptionFunc) (*CustomAttribute, *Response, error) {
	return s.setCustomAttribute("projects", project, c, options...)
}

func (s *CustomAttributesService) setCustomAttribute(resource string, id int64, c CustomAttribute, options ...RequestOptionFunc) (*CustomAttribute, *Response, error) {
	res, resp, err := do[*CustomAttribute](s.client,
		withMethod(http.MethodPut),
		withPath("%s/%d/custom_attributes/%s", resource, id, c.Key),
		withAPIOpts(c),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

func (s *CustomAttributesService) DeleteCustomUserAttribute(user int64, key string, options ...RequestOptionFunc) (*Response, error) {
	return s.deleteCustomAttribute("users", user, key, options...)
}

func (s *CustomAttributesService) DeleteCustomGroupAttribute(group int64, key string, options ...RequestOptionFunc) (*Response, error) {
	return s.deleteCustomAttribute("groups", group, key, options...)
}

func (s *CustomAttributesService) DeleteCustomProjectAttribute(project int64, key string, options ...RequestOptionFunc) (*Response, error) {
	return s.deleteCustomAttribute("projects", project, key, options...)
}

func (s *CustomAttributesService) deleteCustomAttribute(resource string, id int64, key string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("%s/%d/custom_attributes/%s", resource, id, key),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}
