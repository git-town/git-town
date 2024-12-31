@smoke
Feature: show the configuration

  Background:
    Given a Git repo with origin

  Scenario: all configured in Git, no stacked changes
    Given Git setting "git-town.perennial-branches" is "qa staging"
    And Git setting "git-town.perennial-regex" is "^release-"
    And Git setting "git-town.contribution-regex" is "^renovate/"
    And Git setting "git-town.observed-branches" is "observed-1 observed-2"
    And Git setting "git-town.observed-regex" is "^dependabot/"
    And Git setting "git-town.contribution-branches" is "contribution-1 contribution-2"
    And Git setting "git-town.parked-branches" is "parked-1 parked-2"
    And Git setting "git-town.default-branch-type" is "observed"
    And Git setting "git-town.feature-regex" is "^user-.*$"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: ^renovate/
        default branch type: observed
        feature regex: ^user-.*$
        main branch: main
        observed branches: observed-1, observed-2
        observed regex: ^dependabot/
        parked branches: parked-1, parked-2
        perennial branches: qa, staging
        perennial regex: ^release-

      Configuration:
        offline: no

      Create:
        new branch type: feature
        push new branches: no

      Hosting:
        hosting platform: (not set)
        hostname: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)

      Ship:
        delete the tracking branch: yes
        strategy: squash-merge

      Sync:
        run pre-push hook: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync-prototype strategy: merge
        sync tags: yes
        sync with upstream: yes
      """

  Scenario: all configured in config file
    And the configuration file:
      """
      [branches]
      main = "main"
      perennials = [ "public", "staging" ]
      perennial-regex = "^release-"
      default-type = "observed"
      feature-regex = "^user-.*$"
      contribution-regex = "^renovate/"
      observed-regex = "^dependabot/"

      [create]
      push-new-branches = true

      [hosting]
      platform = "github"
      origin-hostname = "github.com"

      [ship]
      delete-tracking-branch = true
      strategy = "squash-merge"

      [sync]
      feature-strategy = "rebase"
      perennial-strategy = "merge"
      prototype-strategy = "compress"
      tags = false
      upstream = true
      """
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: (none)
        contribution regex: ^renovate/
        default branch type: observed
        feature regex: ^user-.*$
        main branch: main
        observed branches: (none)
        observed regex: ^dependabot/
        parked branches: (none)
        perennial branches: public, staging
        perennial regex: ^release-

      Configuration:
        offline: no

      Create:
        new branch type: feature
        push new branches: yes

      Hosting:
        hosting platform: github
        hostname: github.com
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)

      Ship:
        delete the tracking branch: yes
        strategy: squash-merge

      Sync:
        run pre-push hook: yes
        sync-feature strategy: rebase
        sync-perennial strategy: merge
        sync-prototype strategy: compress
        sync tags: no
        sync with upstream: yes
      """

  Scenario: configured in both Git and config file
    Given the main branch is "git-main"
    And Git setting "git-town.perennial-branches" is "git-perennial-1 git-perennial-2"
    And Git setting "git-town.observed-branches" is "observed-1 observed-2"
    And Git setting "git-town.contribution-branches" is "contribution-1 contribution-2"
    And Git setting "git-town.contribution-regex" is "^git-contribution-regex"
    And Git setting "git-town.observed-regex" is "^git-observed-regex"
    And Git setting "git-town.parked-branches" is "parked-1 parked-2"
    And Git setting "git-town.perennial-regex" is "^git-perennial-"
    And Git setting "git-town.feature-regex" is "git-feature-.*"
    And Git setting "git-town.default-branch-type" is "observed"
    And Git setting "git-town.push-new-branches" is "false"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And Git setting "git-town.ship-delete-tracking-branch" is "false"
    And Git setting "git-town.sync-upstream" is "false"
    And Git setting "git-town.sync-tags" is "false"
    And Git setting "git-town.sync-perennial-strategy" is "merge"
    And Git setting "git-town.sync-feature-strategy" is "merge"
    And Git setting "git-town.sync-prototype-strategy" is "compress"
    And the configuration file:
      """
      [branches]
      main = "config-main"
      perennials = [ "config-perennial-1", "config-perennial-2" ]
      perennial-regex = "^config-perennial-"
      default-type = "contribution"
      feature-regex = "^config-feature-.*$"
      contribution-regex = "^config-contribution-regex"
      observed-regex = "^config-observed-regex"

      [create]
      push-new-branches = true

      [hosting]
      platform = "github"
      origin-hostname = "github.com"

      [ship]
      delete-tracking-branch = true
      strategy = "api"

      [sync]
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
        default branch type: observed
        feature regex: git-feature-.*
        main branch: git-main
        observed branches: observed-1, observed-2
        observed regex: ^git-observed-regex
        parked branches: parked-1, parked-2
        perennial branches: git-perennial-1, git-perennial-2, config-perennial-1, config-perennial-2
        perennial regex: ^git-perennial-

      Configuration:
        offline: no

      Create:
        new branch type: feature
        push new branches: no

      Hosting:
        hosting platform: github
        hostname: github.com
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)

      Ship:
        delete the tracking branch: no
        strategy: squash-merge

      Sync:
        run pre-push hook: yes
        sync-feature strategy: merge
        sync-perennial strategy: merge
        sync-prototype strategy: compress
        sync tags: no
        sync with upstream: no
      """

  Scenario: all configured, with stacked changes
    Given the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | alpha  | feature | main   | local     |
      | qa     | (none)  |        | local     |
      | beta   | feature | main   | local     |
      | child  | feature | alpha  | local     |
      | hotfix | feature | qa     | local     |
    And Git setting "git-town.perennial-branches" is "qa"
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: (none)
        contribution regex: (not set)
        default branch type: feature
        feature regex: (not set)
        main branch: main
        observed branches: (none)
        observed regex: (not set)
        parked branches: (none)
        perennial branches: qa
        perennial regex: (not set)

      Configuration:
        offline: no

      Create:
        new branch type: feature
        push new branches: no

      Hosting:
        hosting platform: (not set)
        hostname: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)

      Ship:
        delete the tracking branch: yes
        strategy: api

      Sync:
        run pre-push hook: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync-prototype strategy: merge
        sync tags: yes
        sync with upstream: yes

      Branch Lineage:
        main
          alpha
            child
          beta

        qa
          hotfix
      """

  Scenario: no configuration data
    Given Git Town is not configured
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: (none)
        contribution regex: (not set)
        default branch type: feature
        feature regex: (not set)
        main branch: (not set)
        observed branches: (none)
        observed regex: (not set)
        parked branches: (none)
        perennial branches: (none)
        perennial regex: (not set)

      Configuration:
        offline: no

      Create:
        new branch type: feature
        push new branches: no

      Hosting:
        hosting platform: (not set)
        hostname: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)

      Ship:
        delete the tracking branch: yes
        strategy: api

      Sync:
        run pre-push hook: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync-prototype strategy: merge
        sync tags: yes
        sync with upstream: yes
      """
