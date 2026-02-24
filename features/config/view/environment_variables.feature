Feature: display configuration defined in environment variables

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE         | PARENT | LOCATIONS     |
      | contribution-1 | contribution |        | local, origin |
      | contribution-2 | contribution |        | local, origin |
      | observed-1     | observed     |        | local, origin |
      | observed-2     | observed     |        | local, origin |
    And Git Town is not configured

  Scenario: all configured in env vars, no stacked changes
    When I run "git-town config" with these environment variables
      | GIT_TOWN_AUTO_RESOLVE                | false              |
      | GIT_TOWN_AUTO_SYNC                   | false              |
      | GIT_TOWN_BITBUCKET_APP_PASSWORD      | bitbucket-password |
      | GIT_TOWN_BITBUCKET_USERNAME          | bitbucket-user     |
      | GIT_TOWN_BRANCH_PREFIX               | acme-              |
      | BROWSER                              | firefox            |
      | GIT_TOWN_FORGEJO_TOKEN               | forgejo-token      |
      | GIT_TOWN_CONTRIBUTION_REGEX          | ^renovate/         |
      | GIT_TOWN_DETACHED                    | true               |
      | GIT_TOWN_DEV_REMOTE                  | my-fork            |
      | GIT_TOWN_DISPLAY_TYPES               | all                |
      | GIT_TOWN_FEATURE_REGEX               | ^user-.*$          |
      | GIT_TOWN_FORGE_TYPE                  | gitlab             |
      | GIT_TOWN_GITEA_TOKEN                 | gitea-token        |
      | GIT_TOWN_GITHUB_CONNECTOR            | gh                 |
      | GIT_TOWN_GITHUB_TOKEN                | github-token       |
      | GIT_TOWN_GITLAB_CONNECTOR            | glab               |
      | GIT_TOWN_GITLAB_TOKEN                | gitlab-token       |
      | GIT_TOWN_IGNORE_UNCOMMITTED          | true               |
      | GIT_TOWN_MAIN_BRANCH                 | dev                |
      | GIT_TOWN_NEW_BRANCH_TYPE             | prototype          |
      | GIT_TOWN_OBSERVED_REGEX              | ^dependabot/       |
      | GIT_TOWN_ORDER                       | desc               |
      | GIT_TOWN_PROPOSAL_BREADCRUMB         | stacks             |
      | GIT_TOWN_ORIGIN_HOSTNAME             | codeforge          |
      | GIT_TOWN_OFFLINE                     | 1                  |
      | GIT_TOWN_PERENNIAL_BRANCHES          | qa staging         |
      | GIT_TOWN_PERENNIAL_REGEX             | ^release-          |
      | GIT_TOWN_PUSH_BRANCHES               | no                 |
      | GIT_TOWN_PUSH_HOOK                   | no                 |
      | GIT_TOWN_SHARE_NEW_BRANCHES          | push               |
      | GIT_TOWN_SHIP_DELETE_TRACKING_BRANCH | 0                  |
      | GIT_TOWN_SHIP_STRATEGY               | fast-forward       |
      | GIT_TOWN_STASH                       | false              |
      | GIT_TOWN_SYNC_FEATURE_STRATEGY       | rebase             |
      | GIT_TOWN_SYNC_PERENNIAL_STRATEGY     | merge              |
      | GIT_TOWN_SYNC_PROTOTYPE_STRATEGY     | compress           |
      | GIT_TOWN_SYNC_TAGS                   | false              |
      | GIT_TOWN_SYNC_UPSTREAM               | off                |
      | GIT_TOWN_UNKNOWN_BRANCH_TYPE         | observed           |
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: ^renovate/
        feature regex: ^user-.*$
        main branch: dev
        observed branches: observed-1, observed-2
        observed regex: ^dependabot/
        parked branches: (none)
        perennial branches: qa, staging
        perennial regex: ^release-
        prototype branches: (none)
        unknown branch type: observed
        order: desc
        display types: all branch types

      Configuration:
        offline: yes
        git user name: user
        git user email: email@example.com

      Create:
        branch prefix: acme-
        new branch type: prototype
        share new branches: push
        stash uncommitted changes: no

      Hosting:
        browser: firefox
        development remote: my-fork
        forge type: gitlab
        origin hostname: codeforge
        Bitbucket username: bitbucket-user
        Bitbucket app password: bitbucket-password
        Forgejo token: forgejo-token
        Gitea token: gitea-token
        GitHub connector: gh
        GitHub token: github-token
        GitLab connector: glab
        GitLab token: gitlab-token

      Propose:
        breadcrumb: stacks
        breadcrumb direction: down

      Ship:
        delete tracking branch: no
        ignore uncommitted changes: yes
        ship strategy: fast-forward

      Sync:
        auto-resolve phantom conflicts: no
        auto-sync: no
        run detached: yes
        run pre-push hook: no
        feature sync strategy: rebase
        perennial sync strategy: merge
        prototype sync strategy: compress
        push branches: no
        sync tags: no
        sync with upstream: no
      """
