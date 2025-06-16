## [0.130.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.130.0...v0.130.1) (2025-06-11)


### Bug Fixes

* add missing nil check on create group with avatar ([3298a05](https://gitlab.com/gitlab-org/api/client-go/commit/3298a058f36962a86dea31587956863cd1ed7624))

# [0.130.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.129.0...v0.130.0) (2025-06-11)


### Bug Fixes

* **workflow:** the `release.config.mjs` file mustn't be hidden ([5d423a5](https://gitlab.com/gitlab-org/api/client-go/commit/5d423a55d5b577ebff50dc1a0905c6511b5a4d6f))


### Features

* add "emoji_events" support to group hooks ([c6b770f](https://gitlab.com/gitlab-org/api/client-go/commit/c6b770f350b11e1c9a7c4702ab25b865624b0d47))
* Add `active` to ListProjects ([7818155](https://gitlab.com/gitlab-org/api/client-go/commit/78181558db20647c22e7fed23e749ecafedad27b))
* add generated_file field for MergeRequestDiff ([4b95dac](https://gitlab.com/gitlab-org/api/client-go/commit/4b95dac3ef2b5aabe3040f592ba6378d081d7642))
* add support for `administrator` to Group `project_creation_level` enums ([664bbd7](https://gitlab.com/gitlab-org/api/client-go/commit/664bbd7e3c955c8068b895b1cf1540054ebc13c1))
* add the `WithTokenSource` client option ([6ccfcf8](https://gitlab.com/gitlab-org/api/client-go/commit/6ccfcf857a0a4a850168ecf9317e2e0b8a532173))
* add url field to MergeCommentEvent.merge_request ([bd639d8](https://gitlab.com/gitlab-org/api/client-go/commit/bd639d811c8e7965f426c2deccee84a12d32920f))
* implement a specialized `TokenSource` interface ([83c2e06](https://gitlab.com/gitlab-org/api/client-go/commit/83c2e06cbe76b5268e55589e8bc580582e65bb22))
* **projects:** add ci_push_repository_for_job_token_allowed parameter ([3d539f6](https://gitlab.com/gitlab-org/api/client-go/commit/3d539f66fd63ce4fec6fa7e4e546c9d2acd018f0))
* **terraform-states:** add Terraform States API ([082b81c](https://gitlab.com/gitlab-org/api/client-go/commit/082b81cd456d4b8020f6542daeb3f47c80ba38d0))
