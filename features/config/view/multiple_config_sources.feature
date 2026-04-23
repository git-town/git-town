@smoke
Feature: show the configuration

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE         | PARENT | LOCATIONS     |
      | contribution-1 | contribution |        | local, origin |
      | contribution-2 | contribution |        | local, origin |
      | observed-1     | observed     |        | local, origin |
      | observed-2     | observed     |        | local, origin |
      | parked-1       | parked       | main   | local         |
      | parked-2       | parked       | main   | local         |
      | perennial-1    | perennial    |        | local         |
      | perennial-2    | perennial    |        | local         |
      | prototype-1    | prototype    | main   | local         |
      | prototype-2    | prototype    | main   | local         |

  Scenario: configured in both Git and config file
    Given the main branch is "git-main"
    And Git setting "git-town.auto-sync" is "false"
    And Git setting "git-town.branch-prefix" is "git-"
    And Git setting "git-town.browser" is "firefox"
    And Git setting "git-town.contribution-regex" is "^git-contribution-regex"
    And Git setting "git-town.detached" is "true"
    And Git setting "git-town.display-types" is "no"
    And Git setting "git-town.feature-regex" is "git-feature-.*"
    And Git setting "git-town.github-connector" is "api"
    And Git setting "git-town.ignore-uncommitted" is "false"
    And Git setting "git-town.observed-regex" is "^git-observed-regex"
    And Git setting "git-town.order" is "desc"
    And Git setting "git-town.perennial-branches" is "git-perennial-1 git-perennial-2"
    And Git setting "git-town.perennial-regex" is "^git-perennial-"
    And Git setting "git-town.proposal-breadcrumb" is "stacks"
    And Git setting "git-town.share-new-branches" is "no"
    And Git setting "git-town.ship-delete-tracking-branch" is "false"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And Git setting "git-town.stash" is "false"
    And Git setting "git-town.sync-feature-strategy" is "merge"
    And Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And Git setting "git-town.sync-prototype-strategy" is "compress"
    And Git setting "git-town.sync-tags" is "false"
    And Git setting "git-town.sync-upstream" is "false"
    And Git setting "git-town.unknown-branch-type" is "observed"
    And the configuration file:
      """
      [branches]
      main = "config-main"
      perennials = [ "config-perennial-1", "config-perennial-2" ]
      perennial-regex = "^config-perennial-"
      feature-regex = "^config-feature-.*$"
      contribution-regex = "^config-contribution-regex"
      observed-regex = "^config-observed-regex"
      unknown-type = "contribution"
      order = "asc"
      display-types = "all"

      [create]
      branch-prefix = "file-"
      share-new-branches = "push"
      stash = true

      [hosting]
      browser = "chrome"
      forge-type = "github"
      github-connector = "gh"
      origin-hostname = "github.com"

      [propose]
      breadcrumb = "stacks"

      [ship]
      delete-tracking-branch = true
      ignore-uncommitted = true
      strategy = "api"

      [sync]
      auto-sync = true
      detached = false
      feature-strategy = "merge"
      perennial-strategy = "rebase"
      prototype-strategy = "rebase"
      tags = true
      upstream = true
      """
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: ^git-contribution-regex
        feature regex: git-feature-.*
        main branch: git-main
        observed branches: observed-1, observed-2
        observed regex: ^git-observed-regex
        parked branches: parked-1, parked-2
        perennial branches: git-perennial-1, git-perennial-2, config-perennial-1, config-perennial-2
        perennial regex: ^git-perennial-
        prototype branches: prototype-1, prototype-2
        unknown branch type: observed
        order: desc
        display types: no branch types

      Configuration:
        offline: no
        git user name: user
        git user email: email@example.com

      Create:
        branch prefix: git-
        new branch type: (not set)
        share new branches: no
        stash uncommitted changes: no

      Hosting:
        browser: firefox
        development remote: origin
        forge type: github
        origin hostname: github.com
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Forgejo token: (not set)
        Gitea token: (not set)
        GitHub connector: api
        GitHub token: (not set)
        GitLab connector: (not set)
        GitLab token: (not set)

      Propose:
        breadcrumb: stacks
        breadcrumb direction: down

      Ship:
        delete tracking branch: no
        ignore uncommitted changes: no
        ship strategy: squash-merge

      Sync:
        auto-resolve phantom conflicts: yes
        auto-sync: no
        run detached: yes
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: ff-only
        prototype sync strategy: compress
        push branches: yes
        sync tags: no
        sync with upstream: no
      """
