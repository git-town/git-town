@smoke
Feature: show the configuration

  Scenario: all configured in Git, no stacked changes
    Given the main branch is "main"
    And the perennial branches are "qa" and "staging"
    And local Git Town setting "perennial-regex" is "release-.*"
    And the observed branches "other-1" and "other-2"
    When I run "git-town config"
    Then it prints:
      """
      Branches:
        main branch: main
        perennial branches: qa, staging
        perennial regex: release-.*
        observed branches: other-1, other-2

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship deletes the tracking branch: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync with upstream: yes
        sync before shipping: no

      Hosting:
        hosting platform override: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """

  Scenario: all configured in config file
    Given the configuration file:
      """
      push-new-branches = true
      ship-delete-tracking-branch = true
      sync-upstream = true

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
        observed branches: (none)

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: yes
        ship deletes the tracking branch: yes
        sync-feature strategy: rebase
        sync-perennial strategy: merge
        sync with upstream: yes
        sync before shipping: no

      Hosting:
        hosting platform override: github
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """

  Scenario: configured in both Git and config file
    Given the main branch is "git-main"
    And the perennial branches are "git-perennial-1" and "git-perennial-2"
    And the observed branches "observed-1" and "observed-2"
    And Git Town setting "perennial-regex" is "git-perennial-.*"
    And Git Town setting "push-new-branches" is "false"
    And Git Town setting "ship-delete-tracking-branch" is "false"
    And Git Town setting "sync-upstream" is "false"
    And Git Town setting "sync-perennial-strategy" is "merge"
    And Git Town setting "sync-feature-strategy" is "merge"
    And the configuration file:
      """
      push-new-branches = true
      ship-delete-tracking-branch = true
      sync-upstream = true

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
        perennial branches: config-perennial-1, config-perennial-2, git-perennial-1, git-perennial-2
        perennial regex: git-perennial-.*
        observed branches: observed-1, observed-2

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship deletes the tracking branch: no
        sync-feature strategy: merge
        sync-perennial strategy: merge
        sync with upstream: no
        sync before shipping: no

      Hosting:
        hosting platform override: github
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """

  Scenario: all configured, with stacked changes
    Given the perennial branches "qa" and "staging"
    And the feature branches "alpha" and "beta"
    And a feature branch "child" as a child of "alpha"
    And a feature branch "hotfix" as a child of "qa"
    When I run "git-town config"
    Then it prints:
      """
      Branches:
        main branch: main
        perennial branches: qa, staging
        perennial regex: (not set)
        observed branches: (none)

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship deletes the tracking branch: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync with upstream: yes
        sync before shipping: no

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
        observed branches: (none)

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship deletes the tracking branch: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync with upstream: yes
        sync before shipping: no

      Hosting:
        hosting platform override: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """
