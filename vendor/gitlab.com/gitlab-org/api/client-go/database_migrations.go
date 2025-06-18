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
)

type (
	DatabaseMigrationsServiceInterface interface {
		MarkMigrationAsSuccessful(version int, opt *MarkMigrationAsSuccessfulOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// DatabaseMigrationsService handles communication with the database
	// migrations related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/database_migrations/
	DatabaseMigrationsService struct {
		client *Client
	}
)

var _ DatabaseMigrationsServiceInterface = (*DatabaseMigrationsService)(nil)

// MarkMigrationAsSuccessfulOptions represents the options to mark a migration
// as successful.
//
// GitLab API docs:
// https://docs.gitlab.com/api/database_migrations/#mark-a-migration-as-successful
type MarkMigrationAsSuccessfulOptions struct {
	Database string `url:"database,omitempty" json:"database,omitempty"`
}

// MarkMigrationAsSuccessful markd pending migrations as successfully executed
// to prevent them from being executed by the db:migrate tasks. Use this API to
// skip failing migrations after they are determined to be safe to skip.
//
// GitLab API docs:
// https://docs.gitlab.com/api/database_migrations/#mark-a-migration-as-successful
func (s *DatabaseMigrationsService) MarkMigrationAsSuccessful(version int, opt *MarkMigrationAsSuccessfulOptions, options ...RequestOptionFunc) (*Response, error) {
	u := fmt.Sprintf("admin/migrations/%d/mark", version)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
