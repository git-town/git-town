## 1.46.0

### 🚀 Features

- Update urls and add missing ([!2785](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2785)) by [Mohamed Mongy](https://gitlab.com/mohamedmongy96)



# [1.46.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.45.0...v1.46.0) (2026-03-01)

## 1.45.0

### 🚀 Features

- Add LockMembershipsToSAML support to Application Settings ([!2791](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2791)) by [Mohamed Mongy](https://gitlab.com/mohamedmongy96)

### 🔄 Other Changes

- test(integration): Use epic IID instead of ID in `DeleteEpic` cleanup. ([!2794](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2794)) by [Florian Forster](https://gitlab.com/fforster)



# [1.45.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.44.0...v1.45.0) (2026-02-27)


### Bug Fixes

* **test:** Use epic IID instead of ID in `DeleteEpic` cleanup. ([49dea05](https://gitlab.com/gitlab-org/api/client-go/commit/49dea0587894cd75d5962e69080974fccedde406))

## 1.44.0

### 🚀 Features

- Implement runner controller instance-level runner scope support ([!2765](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2765)) by [Timo Furrer](https://gitlab.com/timofurrer)

### 🔄 Other Changes

- chore(deps): update module github.com/graph-gophers/graphql-go to v1.9.0 ([!2789](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2789)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.44.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.43.0...v1.44.0) (2026-02-26)

## 1.43.0

### 🚀 Features

- feat(pagination): Add `ScanAndCollectN` to collect at most _n_ results. ([!2788](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2788)) by [Florian Forster](https://gitlab.com/fforster)



# [1.43.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.42.0...v1.43.0) (2026-02-25)


### Features

* **pagination:** Add `ScanAndCollectN` to collect at most _n_ results. ([f821c08](https://gitlab.com/gitlab-org/api/client-go/commit/f821c08c2a460755a0ae4db08fa468b54cbb4be1))

## 1.42.0

### 🚀 Features

- feat: Add public_email to CreateUserOptions ([!2787](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2787)) by [Mohamed Othman](https://gitlab.com/mohamed.othman27)



# [1.42.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.41.1...v1.42.0) (2026-02-24)


### Features

* Add public_email to CreateUserOptions ([ab1ec31](https://gitlab.com/gitlab-org/api/client-go/commit/ab1ec3131687de457c8518c60150c254cc56fd83))

## 1.41.1

### 🐛 Bug Fixes

- fix: Fixed a set of endpoints where inputs were escaped and should not be escaped ([!2772](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2772)) by [kilianpaquier](https://gitlab.com/u.kilianpaquier)

### 🔄 Other Changes

- Add `primary_domain` and `pages_primary_domain` to Pages structs ([!2786](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2786)) by [Mohamed Mongy](https://gitlab.com/mohamedmongy96)



## [1.41.1](https://gitlab.com/gitlab-org/api/client-go/compare/v1.41.0...v1.41.1) (2026-02-24)


### Bug Fixes

* Fixed a set of endpoints where inputs were escaped and should not be escaped ([d6d7b17](https://gitlab.com/gitlab-org/api/client-go/commit/d6d7b17f0c4d63c2613ae2aed2ea2901e87c7b8b))

## 1.41.0

### 🚀 Features

- feat: Add missing event toggles to Group Slack integration ([!2784](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2784)) by [Mohamed Othman](https://gitlab.com/mohamed.othman27)

### 🔄 Other Changes

- chore(deps): update module buf.build/go/protovalidate to v1.1.3 ([!2783](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2783)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.41.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.40.1...v1.41.0) (2026-02-22)


### Features

* Add missing event toggles to Group Slack integration ([a4e84a2](https://gitlab.com/gitlab-org/api/client-go/commit/a4e84a27f22083a40f351591c5a851ba19b6a7dc))

## 1.40.1

### 🐛 Bug Fixes

- Add missing group API parameters to Go SDK structs ([!2764](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2764)) by [Mohamed Mongy](https://gitlab.com/mohamedmongy96)



## [1.40.1](https://gitlab.com/gitlab-org/api/client-go/compare/v1.40.0...v1.40.1) (2026-02-21)

## 1.40.0

### 🚀 Features

- feat: Add visibility option to listgroupoptions ([!2775](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2775)) by [Mohamed Othman](https://gitlab.com/mohamed.othman27)
- Add missing parameters to MergeRequestDiff struct ([!2767](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2767)) by [Alekhin Sergey](https://gitlab.com/a.sergey)

### 🔄 Other Changes

- chore(oauth): use go:embed to extract the html ([!2740](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2740)) by [Tomas Vik](https://gitlab.com/viktomas)



# [1.40.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.39.0...v1.40.0) (2026-02-21)


### Features

* Add visibility option to listgroupoptions ([ca08a62](https://gitlab.com/gitlab-org/api/client-go/commit/ca08a62935f8d946dc52e35fcef5528c5950c104))

## 1.39.0

### 🚀 Features

- feat: Add hide backlog and closed list properties to IssueBoards ([!2780](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2780)) by [Mohamed Othman](https://gitlab.com/mohamed.othman27)
- feat(groups): Add provider to AddGroupSAMLLinkOptions ([!2776](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2776)) by [Mohamed Othman](https://gitlab.com/mohamed.othman27)



# [1.39.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.38.0...v1.39.0) (2026-02-19)


### Features

* Add hide backlog and closed list properties to IssueBoards ([a66984e](https://gitlab.com/gitlab-org/api/client-go/commit/a66984ee5934bc55b9618f83d16272b6a4ebf94f))
* **groups:** Add provider to AddGroupSAMLLinkOptions ([bb97c7f](https://gitlab.com/gitlab-org/api/client-go/commit/bb97c7f334ab6cab3eb7153457f14a71b9ff0c55))

## 1.38.0

### 🚀 Features

- feat(events): Add missing parameters for label operations and update documentation links ([!2781](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2781)) by [Mohamed Mongy](https://gitlab.com/mohamedmongy96)
- feat(labels): add missing params and edit links ([!2778](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2778)) by [Mohamed Mongy](https://gitlab.com/mohamedmongy96)

### 🔄 Other Changes

- docs: Fix broken GitLab docs anchors for alert_management API ([!2777](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2777)) by [Mohamed Othman](https://gitlab.com/mohamed.othman27)
- docs: Fix broken documentation links in attestations.go ([!2779](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2779)) by [Mohamed Othman](https://gitlab.com/mohamed.othman27)



# [1.38.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.37.0...v1.38.0) (2026-02-19)


### Features

* **events:** Add missing parameters for label operations and update documentation links ([11b9f08](https://gitlab.com/gitlab-org/api/client-go/commit/11b9f08b37a4c2ada9413259282f163f28b94051))
* **labels:** add missing params and edit links ([ec1b92b](https://gitlab.com/gitlab-org/api/client-go/commit/ec1b92bff403c10446ab1ff6566a3a638871bb7e))

## 1.37.0

### 🚀 Features

- Support system & system_action fields for merge event attributes ([!2737](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2737)) by [Artem Mikheev](https://gitlab.com/renbou)

### 🔄 Other Changes

- Update links of geo_sites.go ([!2782](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2782)) by [Mohamed Mongy](https://gitlab.com/mohamedmongy96)
- chore(deps): update dependency golangci-lint to v2.10.1 ([!2770](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2770)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update golangci/golangci-lint docker tag to v2.10.1 ([!2771](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2771)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update golangci/golangci-lint docker tag to v2.10.0 ([!2769](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2769)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update dependency golangci-lint to v2.10.0 ([!2768](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2768)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.37.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.36.0...v1.37.0) (2026-02-19)

## 1.36.0

### 🚀 Features

- feat: add support for google chat APIs ([!2766](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2766)) by [Zubeen](https://gitlab.com/syedzubeen)

### 🔄 Other Changes

- chore(deps): update module buf.build/go/protovalidate to v1.1.2 ([!2757](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2757)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.36.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.35.0...v1.36.0) (2026-02-17)


### Features

* add support for google chat APIs ([81e58cb](https://gitlab.com/gitlab-org/api/client-go/commit/81e58cbc5296f1ed7651498de367ee42f1a46b1f))

## 1.35.0

### 🚀 Features

- feat(groups): add code_owner_approval_required in a group's default_branch_protection_defaults ([!2725](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2725)) by [Preethi Atchudan](https://gitlab.com/preethiatchudan)

### 🐛 Bug Fixes

- fix(integration): Add missing json tags to ms teams struct ([!2703](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2703)) by [aishahsofea](https://gitlab.com/aishahsofea)

### 🔄 Other Changes

- chore(deps): update module buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go to v1.36.11-20260209202127-80ab13bee0bf.1 ([!2749](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2749)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update node docker tag to v25 ([!2762](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2762)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.35.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.34.0...v1.35.0) (2026-02-16)


### Bug Fixes

* **integration:** Add missing json tags to ms teams struct ([dafd6fd](https://gitlab.com/gitlab-org/api/client-go/commit/dafd6fd9937246278d151e0858aa6cd2a0e8343a))

## 1.34.0

### 🚀 Features

- feat(workitems): Add an initial "Work Items" service with "Get" and "List" methods. ([!2719](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2719)) by [Florian Forster](https://gitlab.com/fforster)

### 🔄 Other Changes

- refactor: migrate to math/rand/v2 ([!2759](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2759)) by [Ville Skyttä](https://gitlab.com/scop)



# [1.34.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.33.0...v1.34.0) (2026-02-13)


### Bug Fixes

* **workitems:** Use `int64` for global work item IDs. ([f04e3d0](https://gitlab.com/gitlab-org/api/client-go/commit/f04e3d08a0e73f535f8049bab43b25753f62cbc0))


### Features

* **request_options:** Add boolean return value to `WithNext`. ([1cd1e1e](https://gitlab.com/gitlab-org/api/client-go/commit/1cd1e1e5ca3ad9c330ada3cbac4f48f22eab9e92))
* **workitems:** Add comprehensive filtering to `ListWorkItemsOptions` ([052a897](https://gitlab.com/gitlab-org/api/client-go/commit/052a897891791acba55afb2fdc5e686ca14ad1df))
* **workitems:** Add pagination support to `ListWorkItems`. ([cfdf5ee](https://gitlab.com/gitlab-org/api/client-go/commit/cfdf5ee61077951a6504b08dfe27033e9bccec5a))
* **workitems:** Add WorkItems service with Get methods ([00925c2](https://gitlab.com/gitlab-org/api/client-go/commit/00925c26114c6a1fb2ad9758ce2ac8658e087f01)), closes [gitlab-org/api/client-go#2213](https://gitlab.com/gitlab-org/api/client-go/issues/2213)
* **workitems:** Implement the `ListWorkItems` method. ([4f8a709](https://gitlab.com/gitlab-org/api/client-go/commit/4f8a7092a23298e3de951564cd0c46a8481c28d7))

## 1.33.0

### 🚀 Features

- Support unauthenticated clients via Unauthenticated auth source ([!2761](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2761)) by [Timo Furrer](https://gitlab.com/timofurrer)



# [1.33.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.32.0...v1.33.0) (2026-02-13)

## 1.32.0

### 🚀 Features

- Implement endpoints for runner controller scopes ([!2758](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2758)) by [Timo Furrer](https://gitlab.com/timofurrer)

### 🔄 Other Changes

- test(namespaces): Address test feedback to simplify the test ([!2744](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2744)) by [Patrick Rice](https://gitlab.com/PatrickRice)
- chore(deps): update golangci/golangci-lint docker tag to v2.9.0 ([!2755](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2755)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update dependency golangci-lint to v2.9.0 ([!2754](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2754)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.32.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.31.0...v1.32.0) (2026-02-13)

## 1.31.0

### 🚀 Features

- Add missing fields to emoji and milestone event types ([!2704](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2704)) by [Huijie Shi](https://gitlab.com/lcdlyxrqy)



# [1.31.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.30.0...v1.31.0) (2026-02-11)

## 1.30.0

### 🚀 Features

- Add missing query params to ListGroupsOptions ([!2726](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2726)) by [Mohamed Mongy](https://gitlab.com/mohamedmongy96)

### 🔄 Other Changes

- chore(deps): update module buf.build/go/protovalidate to v1.1.1 ([!2750](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2750)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- docs(no-release): update url for community fork ([!2748](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2748)) by [Oleksandr Redko](https://gitlab.com/alexandear)



# [1.30.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.29.0...v1.30.0) (2026-02-10)

## 1.29.0

### 🚀 Features

- Update runner controllers to match latest state ([!2747](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2747)) by [Timo Furrer](https://gitlab.com/timofurrer)

### 🔄 Other Changes

- chore(deps): migrate from gopkg.in/yaml.v3 to go.yaml.in/yaml/v3 ([!2639](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2639)) by [Ville Skyttä](https://gitlab.com/scop)



# [1.29.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.28.1...v1.29.0) (2026-02-09)

## 1.28.1

### 🐛 Bug Fixes

- Fix error where GetNamespace double escaped URL-encoded projects ([!2743](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2743)) by [Patrick Rice](https://gitlab.com/PatrickRice)

### 🔄 Other Changes

- refactor: moved comments to interface ([!2716](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2716)) by [Zubeen](https://gitlab.com/syedzubeen)



## [1.28.1](https://gitlab.com/gitlab-org/api/client-go/compare/v1.28.0...v1.28.1) (2026-02-06)

## 1.28.0

### 🚀 Features

- Add `destroy` attribute for pipeline schedule inputs ([!2702](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2702)) by [long nguyen huy](https://gitlab.com/n.h.long.9697)

### 🔄 Other Changes

- Migrate bytes endpoints to new `do` pattern ([!2738](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2738)) by [Timo Furrer](https://gitlab.com/timofurrer)
- docs(users): document the `Locked` and `State` fields ([!2741](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2741)) by [Florian Forster](https://gitlab.com/fforster)
- ci: migrate to Danger from `common-ci-tasks` ([!2742](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2742)) by [Florian Forster](https://gitlab.com/fforster)
- chore(oauth): improve the look of the OAuth confirmation page ([!2739](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2739)) by [Tomas Vik](https://gitlab.com/viktomas)



# [1.28.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.27.0...v1.28.0) (2026-02-05)

## 1.27.0

### 🚀 Features

- Trim leading `@` in user ids in `do()` requests paths ([!2736](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2736)) by [Timo Furrer](https://gitlab.com/timofurrer)

### 🔄 Other Changes

- Migrate endpoints with special status code handling to new `do` pattern ([!2733](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2733)) by [Timo Furrer](https://gitlab.com/timofurrer)
- Support file uploads in `do()` request handler ([!2732](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2732)) by [Timo Furrer](https://gitlab.com/timofurrer)
- Migrate more endpoints to the `do()` pattern ([!2731](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2731)) by [Timo Furrer](https://gitlab.com/timofurrer)
- Revert "refactor(http): preserve response body without copying in multiple services" ([!2730](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2730)) by [Timo Furrer](https://gitlab.com/timofurrer)
- chore(deps): update docker docker tag to v29.2.1 ([!2729](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2729)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.27.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.26.0...v1.27.0) (2026-02-04)

## 1.26.0

### 🚀 Features

- Add slack integration support ([!2692](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2692)) by [Hamza Hassanain](https://gitlab.com/HamzaHassanain)

### 🔄 Other Changes

- refactor(no-release): fix minor revive issues ([!2711](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2711)) by [Oleksandr Redko](https://gitlab.com/alexandear)



# [1.26.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.25.0...v1.26.0) (2026-02-03)

## 1.25.0

### 🚀 Features

- feat(hooks): Add webexintegration ([!2707](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2707)) by [Preethi Atchudan](https://gitlab.com/preethiatchudan)

### 🔄 Other Changes

- fix: Fix broken GitLab docs anchors for Wikis API ([!2723](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2723)) by [Mohamed Othman](https://gitlab.com/mohamed.othman27)
- refactor: moved comments to interface 7 ([!2715](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2715)) by [Zubeen](https://gitlab.com/syedzubeen)
- chore(deps): update module github.com/google/cel-go to v0.27.0 ([!2721](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2721)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- refactor: moved comments to interface 1 ([!2706](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2706)) by [Zubeen](https://gitlab.com/syedzubeen)
- Refactor low complexity endpoints to use new `do` request function ([!2718](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2718)) by [Timo Furrer](https://gitlab.com/timofurrer)
- Add some additional test coverage for functions before migrating to `do` ([!2720](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2720)) by [Patrick Rice](https://gitlab.com/PatrickRice)



# [1.25.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.24.0...v1.25.0) (2026-02-01)


### Bug Fixes

* Fix broken GitLab docs anchors for Wikis API ([bdbb5c0](https://gitlab.com/gitlab-org/api/client-go/commit/bdbb5c0e93847846f6f786c93d649bec18db38e4))


### Features

* **hooks:** Add webexintegration ([857ac6a](https://gitlab.com/gitlab-org/api/client-go/commit/857ac6a82ff63a65ae4df221cf8347fed8946f53))

## 1.24.0

### 🚀 Features

- Add assignee_id to issues api ([!2673](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2673)) by [David Schneider](https://gitlab.com/dvob)



# [1.24.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.23.0...v1.24.0) (2026-01-29)

## 1.23.0

### 🚀 Features

- feat: add group protected branches service ([!2685](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2685)) by [Prakash Divy](https://gitlab.com/prakashdivyy)

### 🔄 Other Changes

- chore(no-release): refactor to slog.DiscardHandler ([!2710](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2710)) by [Oleksandr Redko](https://gitlab.com/alexandear)



# [1.23.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.22.1...v1.23.0) (2026-01-28)


### Features

* add group protected branches service ([c7ffe6f](https://gitlab.com/gitlab-org/api/client-go/commit/c7ffe6ff7bc12996ce27df767a706a253a3ce00b))

## 1.22.1

### 🐛 Bug Fixes

- fix: Type Mismatch in UpdateSettingsOptions for SentryEnabled ([!2690](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2690)) by [Zubeen](https://gitlab.com/syedzubeen)

### 🔄 Other Changes

- fix: URL tags for throttle protected path settings in UpdateSettingsOptions ([!2705](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2705)) by [Zubeen](https://gitlab.com/syedzubeen)
- test: Refactor tests to use testify assertions 7 ([!2700](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2700)) by [Zubeen](https://gitlab.com/syedzubeen)
- test: Refactor tests to use testify assertions 4 ([!2696](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2696)) by [Zubeen](https://gitlab.com/syedzubeen)
- test: Refactor tests to use testify assertions 3 ([!2695](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2695)) by [Zubeen](https://gitlab.com/syedzubeen)
- test(no-release): Refactor tests to use testify assertions 2 ([!2694](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2694)) by [Zubeen](https://gitlab.com/syedzubeen)
- test(no-release): Refactor tests to use testify assertions ([!2693](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2693)) by [Zubeen](https://gitlab.com/syedzubeen)
- test(no-release): Refactor tests to use testify assertions 6 ([!2699](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2699)) by [Zubeen](https://gitlab.com/syedzubeen)



## [1.22.1](https://gitlab.com/gitlab-org/api/client-go/compare/v1.22.0...v1.22.1) (2026-01-28)


### Bug Fixes

* Type Mismatch in UpdateSettingsOptions for SentryEnabled ([c2d3ca9](https://gitlab.com/gitlab-org/api/client-go/commit/c2d3ca98450719f615a951930153ad9fc2585b19))
* URL tags for throttle protected path settings in UpdateSettingsOptions ([a4a525d](https://gitlab.com/gitlab-org/api/client-go/commit/a4a525dce32ba6aa80f45b48fbc0261e59cdabd3))

## 1.22.0

### 🚀 Features

- feat(project_mirror): add ForceSyncProjectMirror ([!2683](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2683)) by [Prakash Divy](https://gitlab.com/prakashdivyy)

### 🔄 Other Changes

- test: Refactor tests to use testify assertions 5 ([!2697](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2697)) by [Zubeen](https://gitlab.com/syedzubeen)



# [1.22.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.21.0...v1.22.0) (2026-01-28)


### Features

* **project_mirror:** add ForceSyncProjectMirror ([b13fcb7](https://gitlab.com/gitlab-org/api/client-go/commit/b13fcb79e6ffb454dc9fd7e332bde90c79a62376))

## 1.21.0

### 🚀 Features

- feat(settings): Add AnonymousSearchesAllowed field support ([!2678](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2678)) by [Seif Hatem](https://gitlab.com/seif-hatem)

### 🔄 Other Changes

- feat: improve URL validation and error handling in client initialization ([!2656](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2656)) by [Amer Khaled](https://gitlab.com/amrkhald777)



# [1.21.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.20.0...v1.21.0) (2026-01-27)


### Features

* improve URL validation and error handling in client initialization ([9417155](https://gitlab.com/gitlab-org/api/client-go/commit/9417155f9c8a5d7c044d052e61d8da5c91bbe57d))
* **settings:** Add AnonymousSearchesAllowed field support ([7185888](https://gitlab.com/gitlab-org/api/client-go/commit/7185888208173e18216ecb353fdfebe91423f0c4))

## 1.20.0

### 🚀 Features

- feat: update events ([!2689](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2689)) by [Huijie Shi](https://gitlab.com/lcdlyxrqy)

### 🔄 Other Changes

- chore(no-release): remove redundant build tag ([!2701](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2701)) by [Oleksandr Redko](https://gitlab.com/alexandear)
- chore(deps): update docker docker tag to v29.2.0 ([!2698](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2698)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.20.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.19.0...v1.20.0) (2026-01-27)


### Features

* update events ([46ba91c](https://gitlab.com/gitlab-org/api/client-go/commit/46ba91cabfe7c13cf4f80738d48ca60b810f520a))

## 1.19.0

### 🚀 Features

- feat(integrations): Add Chat & Notify integrations ([!2691](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2691)) by [Hamza Hassanain](https://gitlab.com/HamzaHassanain)

### 🔄 Other Changes

- refactor: use do function for requests ([!2674](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2674)) by [Timo Furrer](https://gitlab.com/timofurrer)
- chore(docs): Update adding API support guide for new coding patterns ([!2688](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2688)) by [Heidi Berry](https://gitlab.com/heidi.berry)



# [1.19.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.18.0...v1.19.0) (2026-01-26)


### Features

* **integrations:** Add Chat & Notify integrations ([cc692ed](https://gitlab.com/gitlab-org/api/client-go/commit/cc692edd6d8dfed55fd411559af7e53b55d4e2dd))
* **mocks:** add streaming methods for various services ([889b407](https://gitlab.com/gitlab-org/api/client-go/commit/889b407e48432b32b4c1589102ceed6fadb857db))

## 1.18.0

### 🚀 Features

- feat(settings): Added  support for inactive_resource_access_tokens_delete_after_days  to the... ([!2686](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2686)) by [Preethi Atchudan](https://gitlab.com/preethiatchudan)

### 🔄 Other Changes

- Add missing tests for refactored functions ([!2676](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2676)) by [Patrick Rice](https://gitlab.com/PatrickRice)



# [1.18.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.17.0...v1.18.0) (2026-01-25)


### Features

* **settings:** Added  support for inactive_resource_access_tokens_delete_after_days  to the... ([52b60c3](https://gitlab.com/gitlab-org/api/client-go/commit/52b60c3345ef56cc18ae7e8e1e2dd7c9f7f71344))

## 1.17.0

### 🚀 Features

- Add support for Group Mattermost integrations ([!2675](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2675)) by [Hamza Hassanain](https://gitlab.com/HamzaHassanain)



# [1.17.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.16.0...v1.17.0) (2026-01-23)

## 1.16.0

### 🚀 Features

- Add environment, deployed_after, and deployed_before params to merge requests options struct ([!2672](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2672)) by [Filip Aleksic](https://gitlab.com/faleksic)

### 🔄 Other Changes

- chore(deps): update module golang.org/x/oauth2 to v0.34.0 ([!2640](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2640)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.16.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.15.0...v1.16.0) (2026-01-20)

## 1.15.0

### 🚀 Features

- Add ExpiresAt field to ProjectSharedWithGroup struct ([!2671](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2671)) by [cindy](https://gitlab.com/wscix)

### 🔄 Other Changes

- feat: convert examples to testable examples for pkg.go.dev ([!2655](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2655)) by [Amer Khaled](https://gitlab.com/amrkhald777)
- refactor(no-release): enable usetesting linter ([!2664](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2664)) by [Oleksandr Redko](https://gitlab.com/alexandear)
- chore(deps): update docker docker tag to v29.1.5 ([!2665](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2665)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- Draft: Users: Fix GetUsersOptions naming inconsistency ([!2667](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2667)) by [Seif Hatem](https://gitlab.com/seif-hatem)



# [1.15.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.14.0...v1.15.0) (2026-01-20)


### Bug Fixes

* Deprecate incorrect fields returned in response for Emails on Push integration ([71d747d](https://gitlab.com/gitlab-org/api/client-go/commit/71d747da9a297451911b0c4eb4850632a588e3e8))


### Features

* convert examples to testable examples for pkg.go.dev ([fee39f1](https://gitlab.com/gitlab-org/api/client-go/commit/fee39f1f21b264765bbbed80ba23265bd3f633a9))
* **issue_links:** Add ID field to IssueLink struct ([8f813a8](https://gitlab.com/gitlab-org/api/client-go/commit/8f813a8a2e73c41bc81403aceb82d7d94e9ff684))

## 1.14.0

### 🚀 Features

- feat(hooks): Add project hook support for vulnerability events and branch filter strategy ([!2658](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2658)) by [Heidi Berry](https://gitlab.com/heidi.berry)
- Add max_artifacts_size parameter to groups and projects ([!2652](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2652)) by [Betty Godier](https://gitlab.com/betty-godier)



# [1.14.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.13.0...v1.14.0) (2026-01-13)


### Features

* **hooks:** Add project hook support for vulnerability events and branch filter strategy ([4f6d252](https://gitlab.com/gitlab-org/api/client-go/commit/4f6d252a47411602ac6757400e6b5479d807cdb8))

## 1.13.0

### 🚀 Features

- feat(groups): add Active parameter to ListGroupProjects ([!2657](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2657)) by [Kai Armstrong](https://gitlab.com/phikai)

### 🔄 Other Changes

- chore(deps): update docker docker tag to v29.1.4 ([!2651](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2651)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.13.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.12.0...v1.13.0) (2026-01-12)


### Features

* **groups:** add Active parameter to ListGroupProjects ([dec511a](https://gitlab.com/gitlab-org/api/client-go/commit/dec511a199b0adb7ba87f5a02a50651049b68b71))

## 1.12.0

### 🚀 Features

- feat: add EmojiEvents field support to Project Webhooks ([!2653](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2653)) by [Yugan](https://gitlab.com/yugannkt)

### 🔄 Other Changes

- chore(deps): update dependency golangci-lint to v2.8.0 ([!2650](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2650)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- refactor(no-release): use errors.New instead of fmt.Errorf ([!2644](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2644)) by [Oleksandr Redko](https://gitlab.com/alexandear)



# [1.12.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.11.0...v1.12.0) (2026-01-11)


### Features

* add EmojiEvents field support to Project Webhooks ([2bcfa1f](https://gitlab.com/gitlab-org/api/client-go/commit/2bcfa1fd77756a3ccdb2bcf685736ee839b745be))

## 1.11.0

### 🚀 Features

- feat(groups): add support for merge related settings ([!2625](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2625)) by [Zubeen](https://gitlab.com/syedzubeen)

### 🐛 Bug Fixes

- fix(api): typo in ms teams slug ([!2643](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2643)) by [aishahsofea](https://gitlab.com/aishahsofea)

### 🔄 Other Changes

- chore(deps): update module golang.org/x/text to v0.32.0 ([!2642](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2642)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update module golang.org/x/sys to v0.39.0 ([!2641](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2641)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update module github.com/godbus/dbus/v5 to v5.2.2 ([!2637](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2637)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update module github.com/google/go-querystring to v1.2.0 ([!2638](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2638)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.11.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.10.0...v1.11.0) (2026-01-05)


### Bug Fixes

* **api:** typo in ms teams slug ([1ed6c95](https://gitlab.com/gitlab-org/api/client-go/commit/1ed6c9509b23db53c3988a2dde2f11d22c8be5f9))


### Features

* **groups:** add support for merge related settings ([cb8412f](https://gitlab.com/gitlab-org/api/client-go/commit/cb8412fc495d19ee6e44819a2f69fd213d19a199))

## 1.10.0

### 🚀 Features

- feat: implement Runner Controller API ([!2634](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2634)) by [Duo Developer](https://gitlab.com/duo-developer)

### 🔄 Other Changes

- chore(deps): update module github.com/godbus/dbus/v5 to v5.2.1 ([!2635](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2635)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.10.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.9.1...v1.10.0) (2025-12-19)


### Features

* implement Runner Controller API ([66f19f4](https://gitlab.com/gitlab-org/api/client-go/commit/66f19f4073ce87566c7751e0987f857eeb008849))

## 1.9.1

### 🐛 Bug Fixes

- fix: use parameters in config.NewClient and Jobs.DownloadArtifactsFile ([!2633](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2633)) by [Oleksandr Redko](https://gitlab.com/alexandear)

### 🔄 Other Changes

- test: fix TestCreateMergeRequestContextCommits failing locally ([!2631](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2631)) by [Oleksandr Redko](https://gitlab.com/alexandear)
- Code Refactor Using Request Handlers - 8 ([!2523](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2523)) by [Yashesvinee V](https://gitlab.com/yashes7516)
- Code Refactor Using Request Handlers - 6 ([!2521](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2521)) by [Yashesvinee V](https://gitlab.com/yashes7516)



## [1.9.1](https://gitlab.com/gitlab-org/api/client-go/compare/v1.9.0...v1.9.1) (2025-12-17)


### Bug Fixes

* use parameters in config.NewClient and Jobs.DownloadArtifactsFile ([28b7cd7](https://gitlab.com/gitlab-org/api/client-go/commit/28b7cd72f06777a2d3ec7772870c26565140341a))

## 1.9.0

### 🚀 Features

- feat(api): add support for matrix project integration ([!2630](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2630)) by [aishahsofea](https://gitlab.com/aishahsofea)



# [1.9.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.8.2...v1.9.0) (2025-12-16)


### Features

* **api:** add support for matrix project integration ([0a5b11b](https://gitlab.com/gitlab-org/api/client-go/commit/0a5b11b9e2e405fb0a22009d60ce38091cc96625))

## 1.8.2

### 🐛 Bug Fixes

- fix: correct omitempty tag in VariableFilter.EnvironmentScope field ([!2629](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2629)) by [Oleksandr Redko](https://gitlab.com/alexandear)

### 🔄 Other Changes

- feat(protectedTags): add support for `deploy_key_id` to `protected_tags` ([!2624](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2624)) by [Zubeen](https://gitlab.com/syedzubeen)
- chore(deps): update docker docker tag to v29.1.3 ([!2623](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2623)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update module buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go to v1.36.11-20251209175733-2a1774d88802.1 ([!2622](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2622)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update module google.golang.org/protobuf to v1.36.11 ([!2621](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2621)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



## [1.8.2](https://gitlab.com/gitlab-org/api/client-go/compare/v1.8.1...v1.8.2) (2025-12-15)


### Bug Fixes

* correct omitempty tag in VariableFilter.EnvironmentScope field ([c117da1](https://gitlab.com/gitlab-org/api/client-go/commit/c117da1b123251ba86271d1ce3bf9750617e344f))


### Features

* **protectedTags:** add support for `deploy_key_id` to `protected_tags` ([c0fc3db](https://gitlab.com/gitlab-org/api/client-go/commit/c0fc3db793b51bfabb0ac8bb42442e6916b9df3f))

## 1.8.1

### 🐛 Bug Fixes

- fix(epics): handle datetime format in ISOTime UnmarshalJSON ([!2612](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2612)) by [Zubeen](https://gitlab.com/syedzubeen)

### 🔄 Other Changes

- chore(deps): update module buf.build/go/protovalidate to v1.1.0 ([!2619](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2619)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore: add deprecation notice for PersonalAccessTokens.RevokePersonalAccessToken ([!2615](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2615)) by [aishahsofea](https://gitlab.com/aishahsofea)
- chore(deps): update golangci/golangci-lint docker tag to v2.7.2 ([!2613](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2613)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): do not use the experimental package ([!2614](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2614)) by [Mikhail Mazurskiy](https://gitlab.com/ash2k)
- test: Replace SkipIfRunningCE with SkipIfNotLicensed ([!2616](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2616)) by [Patrick Rice](https://gitlab.com/PatrickRice)



## [1.8.1](https://gitlab.com/gitlab-org/api/client-go/compare/v1.8.0...v1.8.1) (2025-12-10)


### Bug Fixes

* **epics:** handle datetime format in ISOTime UnmarshalJSON ([257e0ac](https://gitlab.com/gitlab-org/api/client-go/commit/257e0acd29daf887456d924c0063b52ebc2e808f))

## 1.8.0

### 🚀 Features

- feat(hooks): add support for all hook event types ([!2606](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2606)) by [Heidi Berry](https://gitlab.com/heidi.berry)



# [1.8.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.7.0...v1.8.0) (2025-12-08)


### Features

* **hooks:** add support for all hook event types ([c3c9ca2](https://gitlab.com/gitlab-org/api/client-go/commit/c3c9ca275969adffca37908d63e5c70f634d7bbe))

## 1.7.0

### 🚀 Features

- feat(users): Add support for a user to see only one file diff per page ([!2597](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2597)) by [Zubeen](https://gitlab.com/syedzubeen)



# [1.7.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.6.0...v1.7.0) (2025-12-06)


### Features

* **users:** Add support for a user to see only one file diff per page ([e2a9e09](https://gitlab.com/gitlab-org/api/client-go/commit/e2a9e09e79e7949e0b19dcfc97e3b7b533541856))

## 1.6.0

### 🚀 Features

- feat: add admin compliance policy settings API ([!2610](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2610)) by [Hannes Lange](https://gitlab.com/hlange4)

### 🔄 Other Changes

- doc: fix typo ([!2603](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2603)) by [Guilhem Bonnefille](https://gitlab.com/gbonnefille)
- chore(deps): update golangci/golangci-lint docker tag to v2.7.1 ([!2611](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2611)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update docker docker tag to v29.1.2 ([!2609](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2609)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update golangci/golangci-lint docker tag to v2.7.0 ([!2608](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2608)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.6.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.5.0...v1.6.0) (2025-12-05)


### Features

* add admin compliance policy settings API ([5c17773](https://gitlab.com/gitlab-org/api/client-go/commit/5c17773ca94ddece28978c7396bddcc6c65fb6a7))

## 1.5.0

### 🚀 Features

- feat(Project Mirrors): Add missing Mirror attributes when reading or updating Project Mirrors ([!2600](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2600)) by [Patrick Rice](https://gitlab.com/PatrickRice)



# [1.5.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.4.1...v1.5.0) (2025-12-03)


### Features

* **Project Mirrors:** Add missing Mirror attributes when reading or updating Project Mirrors ([a49b32d](https://gitlab.com/gitlab-org/api/client-go/commit/a49b32df59aeae97247d21a83be3fab97da1bbfe))

## 1.4.1

### 🐛 Bug Fixes

- Encode package managers as CSV in query for dependencies list ([!2604](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2604)) by [Timo Furrer](https://gitlab.com/timofurrer)



## [1.4.1](https://gitlab.com/gitlab-org/api/client-go/compare/v1.4.0...v1.4.1) (2025-12-02)

## 1.4.0

### 🚀 Features

- feat(integrations): Add attestations integrations ([!2582](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2582)) by [Sam Roque-Worcel](https://gitlab.com/sroque-worcel)

### 🔄 Other Changes

- chore(deps): update docker docker tag to v29.1.1 ([!2602](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2602)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.4.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.3.1...v1.4.0) (2025-12-02)


### Features

* **integrations:** Add attestations integrations ([4f50db4](https://gitlab.com/gitlab-org/api/client-go/commit/4f50db4acfb19212bfdfc12eb808dbc7ed8d7ad2))

## 1.3.1

### 🐛 Bug Fixes

- fix(merge_requests): Reinstate missing request option ([!2601](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2601)) by [Heidi Berry](https://gitlab.com/heidi.berry)



## [1.3.1](https://gitlab.com/gitlab-org/api/client-go/compare/v1.3.0...v1.3.1) (2025-12-01)


### Bug Fixes

* **merge_requests:** Reinstate missing request option ([f5f912d](https://gitlab.com/gitlab-org/api/client-go/commit/f5f912ddc2dfb1af88de8710bde783f3f7ccd7c2))

## 1.3.0

### 🚀 Features

- feat(credentials): Add support for revoking group PATs, listing/deleting group SSH keys ([!2594](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2594)) by [Heidi Berry](https://gitlab.com/heidi.berry)

### 🔄 Other Changes

- refactor: moved comments to interface ([!2595](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2595)) by [Zubeen](https://gitlab.com/syedzubeen)
- refactor(users): moved comments to interface ([!2596](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2596)) by [Zubeen](https://gitlab.com/syedzubeen)
- refactor: moved comments to interface ([!2599](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2599)) by [Zubeen](https://gitlab.com/syedzubeen)
- Simplify more request functions, introducing NoEscape ([!2592](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2592)) by [Timo Furrer](https://gitlab.com/timofurrer)



# [1.3.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.2.0...v1.3.0) (2025-11-30)


### Features

* **credentials:** Add support for revoking group PATs, listing/deleting group SSH keys ([3439f4f](https://gitlab.com/gitlab-org/api/client-go/commit/3439f4f0345b97dea0abf926ecaac9d3a7eb6769))

## 1.2.0

### 🚀 Features

- feat(credentials): Add support for listing all SaaS enterprise user personal access tokens ([!2593](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2593)) by [Heidi Berry](https://gitlab.com/heidi.berry)

### 🔄 Other Changes

- Code Refactor Using Request Handlers - 10 ([!2525](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2525)) by [Yashesvinee V](https://gitlab.com/yashes7516)



# [1.2.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.1.0...v1.2.0) (2025-11-27)


### Features

* **credentials:** Add support for listing all SaaS enterprise user personal access tokens ([3697779](https://gitlab.com/gitlab-org/api/client-go/commit/369777938e435b043e37460ff1feffedd84b7dd1))

## 1.1.0

### 🚀 Features

- feat(service_account): allow providing email when update a Service Account ([!2589](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2589)) by [kilianpaquier](https://gitlab.com/u.kilianpaquier)

### 🔄 Other Changes

- Bump dependencies ([!2591](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2591)) by [Timo Furrer](https://gitlab.com/timofurrer)
- chore(deps): update docker docker tag to v29 ([!2586](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2586)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [1.1.0](https://gitlab.com/gitlab-org/api/client-go/compare/v1.0.1...v1.1.0) (2025-11-26)


### Features

* **service_account:** allow providing email when update a Service Account ([324d080](https://gitlab.com/gitlab-org/api/client-go/commit/324d0806a5cd8cb6ae7f68381d09cf5e2a31a0cc))

## 1.0.1

### 🐛 Bug Fixes

- fix: fix ReviewerID() and let it accept int64 ([!2587](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2587)) by [Ilya Savitsky](https://gitlab.com/ipsavitsky234)



## [1.0.1](https://gitlab.com/gitlab-org/api/client-go/compare/v1.0.0...v1.0.1) (2025-11-25)


### Bug Fixes

* fix ReviewerID() and let it accept int64 ([6a6d439](https://gitlab.com/gitlab-org/api/client-go/commit/6a6d43952b70191358e7b726eff4f7f24a0f7ff6))

## 1.0.0

### 💥 Breaking Changes

- Release client-go 1.0 ([!2575](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2575)) by [Patrick Rice](https://gitlab.com/PatrickRice)



# [1.0.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.161.1...v1.0.0) (2025-11-24)


* Merge branch 'release-client-1.0' into 'main' ([f06b8c2](https://gitlab.com/gitlab-org/api/client-go/commit/f06b8c2cb4446e2e76a13bbc707c64e22a64d477))


### Bug Fixes

* **issues:** use AssigneeIDValue for ListProjectIssuesOptions.AssigneeID ([1dcb219](https://gitlab.com/gitlab-org/api/client-go/commit/1dcb219c343bc5b5622ff49933199c003a231bd4))


### Features

* **ListOptions:** Update ListOptions to use composition instead of aliasing ([60beef3](https://gitlab.com/gitlab-org/api/client-go/commit/60beef36d0f93a7dc66749f55d98defbc1b3fe28))


### BREAKING CHANGES

* Release 1.0
* **ListOptions:** ListOptions implementation changed from aliasing to composition
Changelog: Improvements

## 0.161.1

### 🐛 Bug Fixes

- fix(users): Fix a bug where error parsing causes user blocking to not function properly ([!2584](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2584)) by [Patrick Rice](https://gitlab.com/PatrickRice)



## [0.161.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.161.0...v0.161.1) (2025-11-24)


### Bug Fixes

* **users:** Fix a bug where error parsing causes user blocking to not function properly ([2ad5506](https://gitlab.com/gitlab-org/api/client-go/commit/2ad55065d624d27d1f539a3c41489989b9a0d036))

## 0.161.0

### 🚀 Features

- fix: return detailed API errors for BlockUser instead of generic LDAP message ([!2581](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2581)) by [Zubeen](https://gitlab.com/syedzubeen)



# [0.161.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.160.2...v0.161.0) (2025-11-24)


### Bug Fixes

* return detailed API errors for BlockUser instead of generic LDAP message ([2ba9fa6](https://gitlab.com/gitlab-org/api/client-go/commit/2ba9fa6995de6cadf0dae1bf600979b73ee471ce))

## 0.160.2

### 🐛 Bug Fixes

- Fix double escaping in paths ([!2583](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2583)) by [Timo Furrer](https://gitlab.com/timofurrer)



## [0.160.2](https://gitlab.com/gitlab-org/api/client-go/compare/v0.160.1...v0.160.2) (2025-11-24)

## 0.160.1

### 🐛 Bug Fixes

- fix: update input field from "key" to "name" in pipeline schedules to prevent an API error ([!2580](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2580)) by [Zubeen](https://gitlab.com/syedzubeen)

### 🔄 Other Changes

- Code Refactor Using Request Handlers - 9 ([!2524](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2524)) by [Yashesvinee V](https://gitlab.com/yashes7516)
- Code Refactor Using Request Handlers - 7 ([!2522](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2522)) by [Yashesvinee V](https://gitlab.com/yashes7516)
- Code Refactor Using Request Handlers - 5 ([!2518](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2518)) by [Yashesvinee V](https://gitlab.com/yashes7516)
- Code Refactor Using Request Handlers - 2 ([!2515](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2515)) by [Yashesvinee V](https://gitlab.com/yashes7516)
- Code Refactor Using Request Handlers - 4 ([!2517](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2517)) by [Yashesvinee V](https://gitlab.com/yashes7516)
- Code Refactor Using Request Handlers - 3 ([!2516](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2516)) by [Yashesvinee V](https://gitlab.com/yashes7516)
- chore(deps): update module github.com/godbus/dbus/v5 to v5.2.0 ([!2576](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2576)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update golangci/golangci-lint docker tag to v2.6.2 ([!2577](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2577)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- Code Refactor Using Request Handlers - 1 ([!2514](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2514)) by [Yashesvinee V](https://gitlab.com/yashes7516)



## [0.160.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.160.0...v0.160.1) (2025-11-19)


### Bug Fixes

* update input field from "key" to "name" in pipeline schedules to prevent an API error ([062133f](https://gitlab.com/gitlab-org/api/client-go/commit/062133f0c24b32ca6ae64a9f7b80fd3fa7e58256))

## 0.160.0

### 🚀 Features

- feat (project_members): Add show_seat_info option to ProjectMembers ([!2572](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2572)) by [Zubeen](https://gitlab.com/syedzubeen)

### 🔄 Other Changes

- refactor: fix modernize lint issues ([!2574](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2574)) by [Oleksandr Redko](https://gitlab.com/alexandear)
- chore(deps): update module cel.dev/expr to v0.25.1 ([!2573](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2573)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- docs(no-release): format examples, update pkg doc url ([!2543](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2543)) by [Oleksandr Redko](https://gitlab.com/alexandear)



# [0.160.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.159.0...v0.160.0) (2025-11-12)

## 0.159.0

### 🚀 Features

- feat(integrations): add group integration API endpoints for Jira ([!2563](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2563)) by [Harsh Rai](https://gitlab.com/harshrai654)

### 🔄 Other Changes

- chore(deps): update golangci/golangci-lint docker tag to v2.6.1 ([!2564](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2564)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [0.159.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.158.0...v0.159.0) (2025-11-04)


### Features

* **integrations:** add group integration API endpoints for Jira ([09e18ee](https://gitlab.com/gitlab-org/api/client-go/commit/09e18ee598bb7805ac8221f6a05426b1785f9011))

## 0.158.0

### 🚀 Features

- Add support to send variables for GraphQL queries ([!2562](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2562)) by [rafasf](https://gitlab.com/rafasf)

### 🔄 Other Changes

- chore(deps): update module cel.dev/expr to v0.25.0 ([!2560](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2560)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(no-release): standardize GitLab name capitalization ([!2551](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2551)) by [Zubeen](https://gitlab.com/syedzubeen)
- chore(deps): update golangci/golangci-lint docker tag to v2.6.0 ([!2558](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2558)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- refactor: moved comments to interface 2 ([!2557](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2557)) by [Zubeen](https://gitlab.com/syedzubeen)
- refactor: moved comments to interface ([!2556](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2556)) by [Zubeen](https://gitlab.com/syedzubeen)
- refactor(test): avoid panic in tests with goroutines ([!2553](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2553)) by [Oleksandr Redko](https://gitlab.com/alexandear)



# [0.158.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.157.1...v0.158.0) (2025-11-03)

## 0.157.1

### 🐛 Bug Fixes

- fix(protected_packages): fix invalid types ([!2554](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2554)) by [Ruwen Schwedewsky](https://gitlab.com/RuwenSchwedewskySinch)

### 🔄 Other Changes

- chore: Update review instructions for mentioning GitLab ([!2552](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2552)) by [Zubeen](https://gitlab.com/syedzubeen)
- Implement do function to reduce boilerplate ([!2550](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2550)) by [Timo Furrer](https://gitlab.com/timofurrer)
- refactor(test): migrate to testify assertions 4 ([!2548](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2548)) by [Zubeen](https://gitlab.com/syedzubeen)
- refactor(test): migrate to testify assertions 2 ([!2546](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2546)) by [Zubeen](https://gitlab.com/syedzubeen)
- refactor(test): migrate to testify assertions ([!2545](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2545)) by [Zubeen](https://gitlab.com/syedzubeen)
- refactor(test): migrate to testify assertions 5 ([!2549](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2549)) by [Zubeen](https://gitlab.com/syedzubeen)
- test: add unit tests for cluster agents and deployments ([!2499](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2499)) by [Zubeen](https://gitlab.com/syedzubeen)
- refactor(test): migrate to testify assertions 3 ([!2547](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2547)) by [Zubeen](https://gitlab.com/syedzubeen)
- Fix: Helper Functions for Code Refactoring ([!2544](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2544)) by [Yashesvinee V](https://gitlab.com/yashes7516)
- test: adds UT for formatPackageURL ([!2527](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2527)) by [Zubeen](https://gitlab.com/syedzubeen)
- test: adds UT for getEpicLinks ([!2526](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2526)) by [Zubeen](https://gitlab.com/syedzubeen)
- test: add test for ApproveOrRejectProjectDeployment ([!2498](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2498)) by [Zubeen](https://gitlab.com/syedzubeen)
- test: adds UTs for packages ([!2529](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2529)) by [Zubeen](https://gitlab.com/syedzubeen)



## [0.157.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.157.0...v0.157.1) (2025-10-28)


### Bug Fixes

* **no-release:** Helper Functions for Code Refactoring ([6feffea](https://gitlab.com/gitlab-org/api/client-go/commit/6feffea6696a8e333fd0811eee8501e58ba743e3))
* **protected_packages:** fix invalid types ([c09943b](https://gitlab.com/gitlab-org/api/client-go/commit/c09943b0dde510dca32a2544a9c0f75f85943d96))

## 0.157.0

### 🚀 Features

- Add merge requests commit api ([!2539](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2539)) by [Ilya Savitsky](https://gitlab.com/ipsavitsky234)

### 🔄 Other Changes

- test: adds missing UTs for notifications ([!2528](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2528)) by [Zubeen](https://gitlab.com/syedzubeen)
- chore: Update review instructions ([!2537](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2537)) by [Patrick Rice](https://gitlab.com/PatrickRice)
- chore(no-release): Fix godoc comments; enable godoclint ([!2535](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2535)) by [Oleksandr Redko](https://gitlab.com/alexandear)



# [0.157.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.156.0...v0.157.0) (2025-10-13)

## 0.156.0

### 🚀 Features

- feat(api): add support for test report summary api ([!2487](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2487)) by [Daniela Filipe Bento](https://gitlab.com/danifbento)



# [0.156.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.155.0...v0.156.0) (2025-10-10)


### Features

* **api:** add support for test report summary api ([8a0c6dd](https://gitlab.com/gitlab-org/api/client-go/commit/8a0c6dde10a4c9c034274a439eaa060dc6e40995))

## 0.155.0

### 🚀 Features

- feat(group_relations_export): Added Group Relations API integration ([!2508](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2508)) by [Jose Gabriel Companioni Benitez](https://gitlab.com/elC0mpa)

### 🔄 Other Changes

- chore: use local protoc plugin with buf ([!2536](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2536)) by [Timo Furrer](https://gitlab.com/timofurrer)
- chore(no-release): Change generated file comment ([!2532](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2532)) by [Oleksandr Redko](https://gitlab.com/alexandear)
- docs(no-release): Fix the comment for EnvVarGitLabContext ([!2533](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2533)) by [Oleksandr Redko](https://gitlab.com/alexandear)
- feat(client_options): Added unit tests ([!2510](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2510)) by [Jose Gabriel Companioni Benitez](https://gitlab.com/elC0mpa)



# [0.155.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.154.0...v0.155.0) (2025-10-09)


### Features

* **client_options:** Added unit tests ([c148031](https://gitlab.com/gitlab-org/api/client-go/commit/c14803189aa47a0cc9e64e9b455b93e6d4c4e4b9))
* **group_relations_export:** Added Group Relations API integration ([956e039](https://gitlab.com/gitlab-org/api/client-go/commit/956e03950d6bc03c56fa1ea4c5d6e06bfd0b264f))

## 0.154.0

### 🚀 Features

- feat(protected_packages): Add api integration ([!2520](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2520)) by [Jose Gabriel Companioni Benitez](https://gitlab.com/elC0mpa)



# [0.154.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.153.0...v0.154.0) (2025-10-08)


### Features

* **protected_packages:** Add api integration ([2de15c7](https://gitlab.com/gitlab-org/api/client-go/commit/2de15c7875e232b0b0b1e5e5bb8e184cd11d0774))

## 0.153.0

### 🚀 Features

- feat(project_Statistics): Added api integration ([!2512](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2512)) by [Jose Gabriel Companioni Benitez](https://gitlab.com/elC0mpa)

### 🔄 Other Changes

- refactor: moved comments to interface ([!2509](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2509)) by [ajey muthiah](https://gitlab.com/ajeymuthiah)
- chore(no-release): Helper Functions for Code Refactoring ([!2503](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2503)) by [Yashesvinee V](https://gitlab.com/yashes7516)
- Add t.Parallel() to all tests and enable linters ([!2513](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2513)) by [Timo Furrer](https://gitlab.com/timofurrer)
- ci: Remove the `commitlint` job. ([!2511](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2511)) by [Florian Forster](https://gitlab.com/fforster)
- refactor: moved comments to interface ([!2507](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2507)) by [ajey muthiah](https://gitlab.com/ajeymuthiah)



# [0.153.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.152.0...v0.153.0) (2025-10-08)


### Features

* **project_Statistics:** Added api integration ([75b5a03](https://gitlab.com/gitlab-org/api/client-go/commit/75b5a03010a39d5353c975a558fda0b6f00cb697))

## 0.152.0

### 🚀 Features

- feat(api): add api support for listing users who starred a project ([!2486](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2486)) by [ajey muthiah](https://gitlab.com/ajeymuthiah)

### 🔄 Other Changes

- chore(no-release): Update Duo Review Instructions ([!2502](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2502)) by [Patrick Rice](https://gitlab.com/PatrickRice)
- feat(model_registry_api): Added api integration ([!2501](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2501)) by [Jose Gabriel Companioni Benitez](https://gitlab.com/elC0mpa)
- feat(no-release): Add AGENTS.md file ([!2479](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2479)) by [Patrick Rice](https://gitlab.com/PatrickRice)
- chore(no-release): Disable dependency scanning on personal forks ([!2500](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2500)) by [Patrick Rice](https://gitlab.com/PatrickRice)



# [0.152.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.151.0...v0.152.0) (2025-10-06)


### Features

* **api:** add api support for listing users who starred a project ([0cdb4ce](https://gitlab.com/gitlab-org/api/client-go/commit/0cdb4ce5399b43e47bf120a90b16d00c022e194c))
* **model_registry_api:** Added api integration ([065dd63](https://gitlab.com/gitlab-org/api/client-go/commit/065dd639bc8bd0f44cab4d92dbe3ea7f134b913f))
* **no-release:** Add AGENTS.md file ([b9febab](https://gitlab.com/gitlab-org/api/client-go/commit/b9febab3181c3f87edd1fd99b5e596f76bc8b7cc))

## 0.151.0

### 🚀 Features

- feat(api): add api support for delete enterprise user ([!2492](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2492)) by [ajey muthiah](https://gitlab.com/ajeymuthiah)

### 🔄 Other Changes

- docs(no-release): Make it easier to find the docs on issues ([!2497](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2497)) by [Heidi Berry](https://gitlab.com/heidi.berry)



# [0.151.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.150.0...v0.151.0) (2025-10-04)


### Features

* **api:** add api support for delete enterprise user ([36ca8ab](https://gitlab.com/gitlab-org/api/client-go/commit/36ca8ab7672c352a073d59dacae3d763d4089abb))

## 0.150.0

### 🚀 Features

- feat: add Project Aliases API support ([!2493](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2493)) by [Yashesvinee V](https://gitlab.com/yashes7516)

### 🔄 Other Changes

- chore(deps): update module buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go to v1.36.10-20250912141014-52f32327d4b0.1 ([!2495](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2495)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- chore(deps): update module github.com/danieljoos/wincred to v1.2.3 ([!2494](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2494)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)



# [0.150.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.149.0...v0.150.0) (2025-10-03)


### Features

* add Project Aliases API support ([4ece88e](https://gitlab.com/gitlab-org/api/client-go/commit/4ece88e6a8cfa0f53e68184b2905d4c2fb6e857a))

## 0.149.0

### 🚀 Features

- feat(no-release): Add dependency scanning ([!2480](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2480)) by [Patrick Rice](https://gitlab.com/PatrickRice)

### 🔄 Other Changes

- ci(semantic-release): migrate to `@gitlab/semantic-release-merge-request-analyzer` ([!2490](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2490)) by [Florian Forster](https://gitlab.com/fforster)
- ci: add the `autolabels` job ([!2489](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2489)) by [Florian Forster](https://gitlab.com/fforster)
- chore(deps): update module google.golang.org/protobuf to v1.36.10 ([!2488](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2488)) by [GitLab Dependency Bot](https://gitlab.com/gitlab-dependency-update-bot)
- refactor(no-release): added tests for delete project hook method ([!2482](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2482)) by [Jose Gabriel Companioni Benitez](https://gitlab.com/elC0mpa)
- docs(no-release): Add guide for adding new APIs and issue templates ([!2478](https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2478)) by [Heidi Berry](https://gitlab.com/heidi.berry)



# [0.149.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.148.1...v0.149.0) (2025-10-02)


### Features

* **no-release:** Add dependency scanning ([8b0ee10](https://gitlab.com/gitlab-org/api/client-go/commit/8b0ee10acb8adceb5d34be2165b7d587b1e42e49))

## [0.148.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.148.0...v0.148.1) (2025-09-26)


### Bug Fixes

* label unmarshaling for `BasicMergeRequest` list operations ([e80c453](https://gitlab.com/gitlab-org/api/client-go/commit/e80c453aa6a5a265ec8748ae3f3f761a70f4470e))

# [0.148.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.147.1...v0.148.0) (2025-09-23)


### Features

* **ResourceGroup:** add `newest_ready_first` to resource group `process_mode` ([fc8f743](https://gitlab.com/gitlab-org/api/client-go/commit/fc8f7431da4ca8594723105473687e8f1378df2b))

## [0.147.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.147.0...v0.147.1) (2025-09-22)


### Bug Fixes

* **client:** use default retry policy from retryablehttp ([2a72511](https://gitlab.com/gitlab-org/api/client-go/commit/2a725113118608712f668b159ca2dab11f4e588e))

# [0.147.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.146.0...v0.147.0) (2025-09-22)


### Features

* **Project:** add resource_group_default_process_mode ([7804faf](https://gitlab.com/gitlab-org/api/client-go/commit/7804fafa18cc15fec8a0886a081bf3311d72eb1f))

# [0.146.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.145.0...v0.146.0) (2025-09-18)


### Features

* **pipelines:** Add compile-time type-safe pipeline inputs support ([4b30e60](https://gitlab.com/gitlab-org/api/client-go/commit/4b30e60260e4f06e7684352693aac49abd748579)), closes [gitlab-org/api/client-go#2154](https://gitlab.com/gitlab-org/api/client-go/issues/2154)
* **PipelinesService:** Add support for pipeline inputs with type validation ([ab3056f](https://gitlab.com/gitlab-org/api/client-go/commit/ab3056f403ec0268e14b312de3f5b51b115ad97a)), closes [gitlab-org/api/client-go#2154](https://gitlab.com/gitlab-org/api/client-go/issues/2154)
* **PipelineTriggersService:** Add support for pipeline inputs to trigger API ([9ad770e](https://gitlab.com/gitlab-org/api/client-go/commit/9ad770e49e59b2a41c665dfc4781f3b56650e813)), closes [gitlab-org/api/client-go#2154](https://gitlab.com/gitlab-org/api/client-go/issues/2154)

# [0.145.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.144.1...v0.145.0) (2025-09-15)


### Features

* Add missing created_by field to ProjectMembers and GroupMembers ([5348e01](https://gitlab.com/gitlab-org/api/client-go/commit/5348e01913c358c53bdd3da46b069713273d6802))

## [0.144.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.144.0...v0.144.1) (2025-09-13)

# [0.144.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.143.3...v0.144.0) (2025-09-12)


### Features

* **client:** add http.RoundTripper Middleware Configuration Option to Client ([88f9d10](https://gitlab.com/gitlab-org/api/client-go/commit/88f9d1055acbd5e060ab13947b856ccc3a03da6f))

## [0.143.3](https://gitlab.com/gitlab-org/api/client-go/compare/v0.143.2...v0.143.3) (2025-09-10)

## [0.143.2](https://gitlab.com/gitlab-org/api/client-go/compare/v0.143.1...v0.143.2) (2025-09-09)

## [0.143.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.143.0...v0.143.1) (2025-09-08)

# [0.143.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.142.6...v0.143.0) (2025-09-08)


### Features

* **users:** Add support for PublicEmail to ListUsers ([74a3b6a](https://gitlab.com/gitlab-org/api/client-go/commit/74a3b6a7dd1340faa70ec1246b5b99394c56f90b))

## [0.142.6](https://gitlab.com/gitlab-org/api/client-go/compare/v0.142.5...v0.142.6) (2025-09-02)

## [0.142.5](https://gitlab.com/gitlab-org/api/client-go/compare/v0.142.4...v0.142.5) (2025-08-30)

## [0.142.4](https://gitlab.com/gitlab-org/api/client-go/compare/v0.142.3...v0.142.4) (2025-08-28)

## [0.142.3](https://gitlab.com/gitlab-org/api/client-go/compare/v0.142.2...v0.142.3) (2025-08-28)

## [0.142.2](https://gitlab.com/gitlab-org/api/client-go/compare/v0.142.1...v0.142.2) (2025-08-26)

## [0.142.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.142.0...v0.142.1) (2025-08-25)

# [0.142.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.141.2...v0.142.0) (2025-08-21)


### Features

* **tokens:** add expiration filters and sorting options to ListPersonalAccessTokens ([0a9f797](https://gitlab.com/gitlab-org/api/client-go/commit/0a9f79790ac87c7f7b8e32e9cdea27fbc613728b))

## [0.141.2](https://gitlab.com/gitlab-org/api/client-go/compare/v0.141.1...v0.141.2) (2025-08-20)

## [0.141.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.141.0...v0.141.1) (2025-08-18)

# [0.141.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.140.0...v0.141.0) (2025-08-18)


### Features

* **config:** support custom headers for instances ([76b0e82](https://gitlab.com/gitlab-org/api/client-go/commit/76b0e82ab57b21b7da915117fb37ac2bf56506e8))

# [0.140.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.139.2...v0.140.0) (2025-08-18)


### Features

* **client:** add support for cookie jars ([4b525e3](https://gitlab.com/gitlab-org/api/client-go/commit/4b525e3f14741176ea8cbf4e7ae988b87455f4d0))

## [0.139.2](https://gitlab.com/gitlab-org/api/client-go/compare/v0.139.1...v0.139.2) (2025-08-14)

## [0.139.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.139.0...v0.139.1) (2025-08-14)

# [0.139.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.138.0...v0.139.0) (2025-08-13)


### Features

* **terraform:** improve Terraform States service ([e08128b](https://gitlab.com/gitlab-org/api/client-go/commit/e08128bf87011455db06dc946e77b2a16ee36948))

# [0.138.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.137.0...v0.138.0) (2025-08-12)


### Bug Fixes

* deprecate ListProjectInvidedGroupOptions due to a typo ([322496a](https://gitlab.com/gitlab-org/api/client-go/commit/322496a8a4c3fd7393b4b2c2b427c42fff243861))
* Update config package name to v1beta1 ([f958e6b](https://gitlab.com/gitlab-org/api/client-go/commit/f958e6bd2935fddf4867d9992908e87288e89c20))


### Features

* add support for field "Created at" for Tags ([f363d57](https://gitlab.com/gitlab-org/api/client-go/commit/f363d57853f2e05c848e88946269c936f0b6bf76))
* **app settings:** Add support for CanCreateOrganization ([1db661d](https://gitlab.com/gitlab-org/api/client-go/commit/1db661de26e0d3a78134c6bd1d31fb24d9a60677))
* **hooks:** Add support for project webhook url variables ([efabed5](https://gitlab.com/gitlab-org/api/client-go/commit/efabed57d83eefe565aa2dbbb943d94212ec6167))
* update datadog integration with new fields and API endpoints ([660ef31](https://gitlab.com/gitlab-org/api/client-go/commit/660ef31daf884bde545cfaa88432ac5ec7e3bfe7))
* update external status checks to return the status check object ([2d78e8c](https://gitlab.com/gitlab-org/api/client-go/commit/2d78e8cc43971c4395c980672de7263c10401900))

# [0.137.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.136.0...v0.137.0) (2025-07-21)


### Features

* **integrations:** add group harbor integration ([220e4cb](https://gitlab.com/gitlab-org/api/client-go/commit/220e4cb524d9303d36384043f29f96f43e4d9387))

# [0.136.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.135.0...v0.136.0) (2025-07-21)


### Features

* **client:** add NewRequestToURL function for calls to absolute URLs ([524b571](https://gitlab.com/gitlab-org/api/client-go/commit/524b571339b7704e0e346a5a64f367265b96125f))
* **projects:** Add support for RestoreProject ([b33e888](https://gitlab.com/gitlab-org/api/client-go/commit/b33e8882ad6611b1ac19222d0abdbfc477846ea1))

# [0.135.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.134.0...v0.135.0) (2025-07-21)


### Features

* **config:** implement extensions API ([257f745](https://gitlab.com/gitlab-org/api/client-go/commit/257f74599727b6a006ba65b1c3efd7ff5fc7b86c))
* **config:** initial push of the ability to use a config file for auth ([575c0cc](https://gitlab.com/gitlab-org/api/client-go/commit/575c0cc6a1de48582ea9b3b19cef021dc3f1397a))
* **integrations:** add group integration for microsoft teams ([da0b1dd](https://gitlab.com/gitlab-org/api/client-go/commit/da0b1dd5b86fd6a41d7c043621611d0687fc628f))
* **merge-requests:** add auto_merge, deprecate old field, for merging a request ([9119eb0](https://gitlab.com/gitlab-org/api/client-go/commit/9119eb0e6662f136e589cdee74aaa410136ca664))

# [0.134.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.133.1...v0.134.0) (2025-07-07)


### Features

* **oauth:** implement OAuth2 helper package ([a44e8eb](https://gitlab.com/gitlab-org/api/client-go/commit/a44e8eb7743ff8d948f396b9849a82a7d7d6d6c4))

## [0.133.1](https://gitlab.com/gitlab-org/api/client-go/compare/v0.133.0...v0.133.1) (2025-07-07)


### Bug Fixes

* deprecate ProjectReposityStorage due to a typo ([38a9652](https://gitlab.com/gitlab-org/api/client-go/commit/38a965279a4c570fd4db4f08503a63c4e7177439))

# [0.133.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.132.0...v0.133.0) (2025-07-03)


### Features

* **testing:** allow to specify client options when creating test client ([9377147](https://gitlab.com/gitlab-org/api/client-go/commit/93771470166ce7c9097328b5e49f75a381c1720b))

# [0.132.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.131.0...v0.132.0) (2025-07-02)


### Bug Fixes

* **no-release:** fix body-max-line-length ([f5d6d05](https://gitlab.com/gitlab-org/api/client-go/commit/f5d6d05d5781cd4fc31fa647ed94d486a1f6fa72))


### Features

* add missing ref_protected property from PushWebhookEventType ([15d0224](https://gitlab.com/gitlab-org/api/client-go/commit/15d0224575e7a5415783466afffe6c6b7aaf5dec))
* add WithUserAgent client option ([3e8b80c](https://gitlab.com/gitlab-org/api/client-go/commit/3e8b80cd40b3d4ad54cb050ebd1b6e11b848869a))
* export various auth sources ([281e408](https://gitlab.com/gitlab-org/api/client-go/commit/281e4083beed2b88b035dddcb562982d4c412143))
* **serviceaccounts:** bring group service accounts in line with API ([a08974f](https://gitlab.com/gitlab-org/api/client-go/commit/a08974f284c043d4039495ed4b8f24ebeb256cdc))
* **serviceaccounts:** bring group service accounts in line with API ([fb582a4](https://gitlab.com/gitlab-org/api/client-go/commit/fb582a4bb523443984851bc1d4b0fb699cfa2a9f))

# [0.131.0](https://gitlab.com/gitlab-org/api/client-go/compare/v0.130.1...v0.131.0) (2025-07-01)


### Features

* add ScanAndCollect for pagination ([cbac9ae](https://gitlab.com/gitlab-org/api/client-go/commit/cbac9aed9bb3c7f8d175585a6d38baa3f2a7fbe1))
* add support for optional query params to get commit statuses ([e1b29ad](https://gitlab.com/gitlab-org/api/client-go/commit/e1b29adfd37db39aae4e1547f336b71d67efcdb8))

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
