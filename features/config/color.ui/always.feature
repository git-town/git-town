@smoke
Feature: show the configuration

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE         | PARENT | LOCATIONS |
      | contribution-1 | contribution |        | local     |
      | contribution-2 | contribution |        | local     |
      | observed-1     | observed     |        | local     |
      | observed-2     | observed     |        | local     |
      | parked-1       | parked       | main   | local     |
      | parked-2       | parked       | main   | local     |
      | perennial-1    | perennial    |        | local     |
      | perennial-2    | perennial    |        | local     |
      | prototype-1    | prototype    | main   | local     |
      | prototype-2    | prototype    | main   | local     |
    And local Git setting "color.ui" is "always"

  Scenario: all configured in Git, no stacked changes
    Given Git setting "git-town.perennial-branches" is "qa staging"
    And Git setting "git-town.perennial-regex" is "^release-"
    And Git setting "git-town.contribution-regex" is "^renovate/"
    And Git setting "git-town.observed-regex" is "^dependabot/"
    And Git setting "git-town.unknown-branch-type" is "observed"
    And Git setting "git-town.feature-regex" is "^user-.*$"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: ^renovate/
        feature regex: ^user-.*$
        main branch: main
        observed branches: observed-1, observed-2
        observed regex: ^dependabot/
        parked branches: parked-1, parked-2
        perennial branches: qa, staging
        perennial regex: ^release-
        prototype branches: prototype-1, prototype-2
        unknown branch type: observed

      Configuration:
        offline: no

      Create:
        new branch type: (not set)
        share new branches: no

      Hosting:
        development remote: origin
        forge type: (not set)
        origin hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Codeberg token: (not set)
        Gitea token: (not set)
        GitHub token: (not set)
        GitLab token: (not set)

      Ship:
        delete tracking branch: yes
        ship strategy: squash-merge

      Sync:
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: rebase
        prototype sync strategy: merge
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
      feature-regex = "^user-.*$"
      contribution-regex = "^renovate/"
      observed-regex = "^dependabot/"
      unknown-type = "observed"

      [create]
      share-new-branches = "push"

      [hosting]
      forge-type = "github"
      origin-hostname = "github.com"

      [ship]
      delete-tracking-branch = true
      strategy = "squash-merge"

      [sync]
      feature-strategy = "rebase"
      perennial-strategy = "ff-only"
      prototype-strategy = "compress"
      tags = false
      upstream = true
      """
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: ^renovate/
        feature regex: ^user-.*$
        main branch: main
        observed branches: observed-1, observed-2
        observed regex: ^dependabot/
        parked branches: parked-1, parked-2
        perennial branches: public, staging
        perennial regex: ^release-
        prototype branches: prototype-1, prototype-2
        unknown branch type: observed

      Configuration:
        offline: no

      Create:
        new branch type: (not set)
        share new branches: push

      Hosting:
        development remote: origin
        forge type: github
        origin hostname: github.com
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Codeberg token: (not set)
        Gitea token: (not set)
        GitHub token: (not set)
        GitLab token: (not set)

      Ship:
        delete tracking branch: yes
        ship strategy: squash-merge

      Sync:
        run pre-push hook: yes
        feature sync strategy: rebase
        perennial sync strategy: ff-only
        prototype sync strategy: compress
        sync tags: no
        sync with upstream: yes
      """

  Scenario: configured in both Git and config file
    Given the main branch is "git-main"
    And Git setting "git-town.perennial-branches" is "git-perennial-1 git-perennial-2"
    And Git setting "git-town.contribution-regex" is "^git-contribution-regex"
    And Git setting "git-town.observed-regex" is "^git-observed-regex"
    And Git setting "git-town.perennial-regex" is "^git-perennial-"
    And Git setting "git-town.feature-regex" is "git-feature-.*"
    And Git setting "git-town.share-new-branches" is "no"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And Git setting "git-town.ship-delete-tracking-branch" is "false"
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

      [hosting]
      forge-type = "github"
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
        feature regex: git-feature-.*
        main branch: git-main
        observed branches: observed-1, observed-2
        observed regex: ^git-observed-regex
        parked branches: parked-1, parked-2
        perennial branches: git-perennial-1, git-perennial-2, config-perennial-1, config-perennial-2
        perennial regex: ^git-perennial-
        prototype branches: prototype-1, prototype-2
        unknown branch type: observed

      Configuration:
        offline: no

      Create:
        new branch type: (not set)
        share new branches: no

      Hosting:
        development remote: origin
        forge type: github
        origin hostname: github.com
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Codeberg token: (not set)
        Gitea token: (not set)
        GitHub token: (not set)
        GitLab token: (not set)

      Ship:
        delete tracking branch: no
        ship strategy: squash-merge

      Sync:
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: ff-only
        prototype sync strategy: compress
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
        contribution branches: contribution-1, contribution-2
        contribution regex: (not set)
        feature regex: (not set)
        main branch: main
        observed branches: observed-1, observed-2
        observed regex: (not set)
        parked branches: parked-1, parked-2
        perennial branches: qa
        perennial regex: (not set)
        prototype branches: prototype-1, prototype-2
        unknown branch type: feature

      Configuration:
        offline: no

      Create:
        new branch type: (not set)
        share new branches: no

      Hosting:
        development remote: origin
        forge type: (not set)
        origin hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Codeberg token: (not set)
        Gitea token: (not set)
        GitHub token: (not set)
        GitLab token: (not set)

      Ship:
        delete tracking branch: yes
        ship strategy: api

      Sync:
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: rebase
        prototype sync strategy: merge
        sync tags: yes
        sync with upstream: yes

      Branch Lineage:
        main
          alpha
            child
          beta
          parked-1
          parked-2
          prototype-1
          prototype-2

        qa
          hotfix
      """

  Scenario: no configuration data
    Given Git Town is not configured
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: (not set)
        feature regex: (not set)
        main branch: (not set)
        observed branches: observed-1, observed-2
        observed regex: (not set)
        parked branches: parked-1, parked-2
        perennial branches: (none)
        perennial regex: (not set)
        prototype branches: prototype-1, prototype-2
        unknown branch type: feature

      Configuration:
        offline: no

      Create:
        new branch type: (not set)
        share new branches: no

      Hosting:
        development remote: origin
        forge type: (not set)
        origin hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Codeberg token: (not set)
        Gitea token: (not set)
        GitHub token: (not set)
        GitLab token: (not set)

      Ship:
        delete tracking branch: yes
        ship strategy: api

      Sync:
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: rebase
        prototype sync strategy: merge
        sync tags: yes
        sync with upstream: yes
      """
