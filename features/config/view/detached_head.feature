Feature: display configuration from Git metadata in detached head state

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | branch    | feature   | main   | local         |
      | observed  | observed  |        | local, origin |
      | parked    | parked    | main   | local         |
      | perennial | perennial |        | local         |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | branch | local    | commit 1 |
      |        | local    | commit 2 |
    And Git setting "git-town.auto-resolve" is "false"
    And Git setting "git-town.auto-sync" is "false"
    And Git setting "git-town.branch-prefix" is "acme-"
    And Git setting "git-town.browser" is "firefox"
    And Git setting "git-town.contribution-regex" is "^renovate/"
    And Git setting "git-town.detached" is "true"
    And Git setting "git-town.display-types" is "all"
    And Git setting "git-town.feature-regex" is "^user-.*$"
    And Git setting "git-town.ignore-uncommitted" is "true"
    And Git setting "git-town.observed-regex" is "^dependabot/"
    And Git setting "git-town.perennial-branches" is "qa staging"
    And Git setting "git-town.perennial-regex" is "^release-"
    And Git setting "git-town.proposal-breadcrumb" is "branches"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And Git setting "git-town.stash" is "false"
    And Git setting "git-town.unknown-branch-type" is "observed"
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town config"

  Scenario: result
    Then Git Town prints:
      """
      Branches:
        contribution branches: (none)
        contribution regex: ^renovate/
        feature regex: ^user-.*$
        main branch: main
        observed branches: observed
        observed regex: ^dependabot/
        parked branches: parked
        perennial branches: qa, staging
        perennial regex: ^release-
        prototype branches: (none)
        unknown branch type: observed
        order: asc
        display types: all branch types

      Configuration:
        offline: no
        git user name: user
        git user email: email@example.com

      Create:
        branch prefix: acme-
        new branch type: (not set)
        share new branches: no
        stash uncommitted changes: no

      Hosting:
        browser: firefox
        development remote: origin
        forge type: (not set)
        origin hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Forgejo token: (not set)
        Gitea token: (not set)
        GitHub connector: (not set)
        GitHub token: (not set)
        GitLab connector: (not set)
        GitLab token: (not set)

      Propose:
        breadcrumb: branches
        breadcrumb direction: down

      Ship:
        delete tracking branch: yes
        ignore uncommitted changes: yes
        ship strategy: squash-merge

      Sync:
        auto-resolve phantom conflicts: no
        auto-sync: no
        run detached: yes
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: rebase
        prototype sync strategy: merge
        push branches: yes
        sync tags: yes
        sync with upstream: yes
        auto-resolve phantom conflicts: no

      Branch Lineage:
        main
          branch
          parked
      """
