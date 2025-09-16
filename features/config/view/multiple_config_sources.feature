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
      | perennial-1    | perennial    |        | local         |
      | perennial-2    | perennial    |        | local         |

  Scenario: configured in both Git and config file
    Given the main branch is "git-main"
    And Git setting "git-town.perennial-branches" is "git-perennial-1 git-perennial-2"
    And Git setting "git-town.contribution-regex" is "^git-contribution-regex"
    And Git setting "git-town.detached" is "true"
    And Git setting "git-town.observed-regex" is "^git-observed-regex"
    And Git setting "git-town.perennial-regex" is "^git-perennial-"
    And Git setting "git-town.feature-regex" is "git-feature-.*"
    And Git setting "git-town.github-connector" is "api"
    And Git setting "git-town.share-new-branches" is "no"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And Git setting "git-town.ship-delete-tracking-branch" is "false"
    And Git setting "git-town.stash" is "false"
    And Git setting "git-town.sync-upstream" is "false"
    And Git setting "git-town.sync-tags" is "false"
    And Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And Git setting "git-town.sync-feature-strategy" is "merge"
    And Git setting "git-town.sync-prototype-strategy" is "compress"
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

      [create]
      share-new-branches = "push"
      stash = true

      [hosting]
      forge-type = "github"
      github-connector = "gh"
      origin-hostname = "github.com"

      [ship]
      delete-tracking-branch = true
      strategy = "api"

      [sync]
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
        parked branches: (none)
        perennial branches: git-perennial-1, git-perennial-2, config-perennial-1, config-perennial-2
        perennial regex: ^git-perennial-
        prototype branches: (none)
        unknown branch type: observed

      Configuration:
        offline: no

      Create:
        new branch type: (not set)
        share new branches: no
        stash uncommitted changes: no

      Hosting:
        development remote: origin
        forge type: github
        origin hostname: github.com
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Forgejo token: (not set)
        Gitea token: (not set)
        GitHub connector type: api
        GitHub token: (not set)
        GitLab connector type: (not set)
        GitLab token: (not set)

      Ship:
        delete tracking branch: no
        ship strategy: squash-merge

      Sync:
        auto-resolve phantom conflicts: yes
        run detached: yes
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: ff-only
        prototype sync strategy: compress
        push branches: yes
        sync tags: no
        sync with upstream: no
      """
