@smoke
Feature: show the configuration

  Background:
    Given a Git repo with origin

  Scenario: all configured in Git, no stacked changes
    Given the branches
      | NAME           | TYPE         | PARENT | LOCATIONS |
      | feature        | feature      | main   | local     |
      | qa             | perennial    |        | local     |
      | staging        | perennial    |        | local     |
      | observed-1     | observed     |        | local     |
      | observed-2     | observed     |        | local     |
      | contribution-1 | contribution |        | local     |
      | contribution-2 | contribution |        | local     |
      | parked-1       | parked       | main   | local     |
      | parked-2       | parked       | main   | local     |
    And the main branch is "main"
    And local Git Town setting "perennial-regex" is "release-.*"
    When I run "git-town config"
    Then it prints:
      """
      Branches:
        main branch: main
        perennial branches: qa, staging
        perennial regex: release-.*
        parked branches: parked-1, parked-2
        contribution branches: contribution-1, contribution-2
        observed branches: observed-1, observed-2

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship deletes the tracking branch: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync with upstream: yes
        sync tags: yes

      Hosting:
        hosting platform override: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """

  Scenario: all configured in config file
    And the configuration file:
      """
      push-new-branches = true
      ship-delete-tracking-branch = true
      sync-upstream = true
      sync-tags = false

      [branches]
      main = "main"
      perennials = [ "public", "staging" ]
      perennial-regex = "release-.*"

      [hosting]
      platform = "github"
      origin-hostname = "github.com"

      [sync-strategy]
      feature-branches = "rebase"
      perennial-branches = "merge"
      """
    When I run "git-town config"
    Then it prints:
      """
      Branches:
        main branch: main
        perennial branches: public, staging
        perennial regex: release-.*
        parked branches: (none)
        contribution branches: (none)
        observed branches: (none)

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: yes
        ship deletes the tracking branch: yes
        sync-feature strategy: rebase
        sync-perennial strategy: merge
        sync with upstream: yes
        sync tags: no

      Hosting:
        hosting platform override: github
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """

  Scenario: configured in both Git and config file
    Given the branches
      | NAME            | TYPE         | PARENT | LOCATIONS |
      | git-perennial-1 | perennial    |        | local     |
      | git-perennial-2 | perennial    |        | local     |
      | observed-1      | observed     |        | local     |
      | observed-2      | observed     |        | local     |
      | contribution-1  | contribution |        | local     |
      | contribution-2  | contribution |        | local     |
      | parked-1        | parked       | main   | local     |
      | parked-2        | parked       | main   | local     |
    And the main branch is "git-main"
    And Git Town setting "perennial-regex" is "git-perennial-.*"
    And Git Town setting "push-new-branches" is "false"
    And Git Town setting "ship-delete-tracking-branch" is "false"
    And Git Town setting "sync-upstream" is "false"
    And Git Town setting "sync-tags" is "false"
    And Git Town setting "sync-perennial-strategy" is "merge"
    And Git Town setting "sync-feature-strategy" is "merge"
    And the configuration file:
      """
      push-new-branches = true
      ship-delete-tracking-branch = true
      sync-upstream = true
      sync-tags = true

      [branches]
      main = "config-main"
      perennials = [ "config-perennial-1", "config-perennial-2" ]
      perennial-regex = "config-perennial-.*"

      [hosting]
      platform = "github"
      origin-hostname = "github.com"

      [sync-strategy]
      feature-branches = "merge"
      perennial-branches = "merge"
      """
    When I run "git-town config"
    Then it prints:
      """
      Branches:
        main branch: git-main
        perennial branches: git-perennial-1, git-perennial-2, config-perennial-1, config-perennial-2
        perennial regex: git-perennial-.*
        parked branches: parked-1, parked-2
        contribution branches: contribution-1, contribution-2
        observed branches: observed-1, observed-2

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship deletes the tracking branch: no
        sync-feature strategy: merge
        sync-perennial strategy: merge
        sync with upstream: no
        sync tags: no

      Hosting:
        hosting platform override: github
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """

  Scenario: all configured, with stacked changes
    Given the branches
      | NAME    | TYPE      | PARENT | LOCATIONS |
      | alpha   | feature   | main   | local     |
      | qa      | perennial |        | local     |
      | staging | perennial |        | local     |
      | beta    | feature   | main   | local     |
      | child   | feature   | alpha  | local     |
      | hotfix  | feature   | qa     | local     |
    When I run "git-town config"
    Then it prints:
      """
      Branches:
        main branch: main
        perennial branches: qa, staging
        perennial regex: (not set)
        parked branches: (none)
        contribution branches: (none)
        observed branches: (none)

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship deletes the tracking branch: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync with upstream: yes
        sync tags: yes

      Hosting:
        hosting platform override: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)

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
    Then it prints:
      """
      Branches:
        main branch: (not set)
        perennial branches: (none)
        perennial regex: (not set)
        parked branches: (none)
        contribution branches: (none)
        observed branches: (none)

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship deletes the tracking branch: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync with upstream: yes
        sync tags: yes

      Hosting:
        hosting platform override: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """
