@smoke
Feature: show the configuration

  Background:
    Given a Git repo with origin

  Scenario: all configured in Git, no stacked changes
    Given Git Town setting "perennial-branches" is "qa staging"
    And Git Town setting "perennial-regex" is "^release-"
    And Git Town setting "contribution-regex" is "^renovate/"
    And Git Town setting "observed-branches" is "observed-1 observed-2"
    And Git Town setting "observed-regex" is "^dependabot/"
    And Git Town setting "contribution-branches" is "contribution-1 contribution-2"
    And Git Town setting "parked-branches" is "parked-1 parked-2"
    And Git Town setting "default-branch-type" is "observed"
    And Git Town setting "feature-regex" is "^user-.*$"
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town config"
    Then it prints:
      """
      Branches:
        main branch: main
        perennial branches: qa, staging
        perennial regex: ^release-
        default branch type: observed
        feature regex: ^user-.*$
        parked branches: parked-1, parked-2
        contribution branches: contribution-1, contribution-2
        contribution regex: ^renovate/
        observed branches: observed-1, observed-2
        observed regex: ^dependabot/

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship strategy: squash-merge
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

  @this
  Scenario: all configured in config file
    And the configuration file:
      """
      push-new-branches = true
      ship-strategy = "squash-merge"
      ship-delete-tracking-branch = true
      sync-upstream = true
      sync-tags = false

      [branches]
      main = "main"
      perennials = [ "public", "staging" ]
      perennial-regex = "^release-"
      default-type = "observed"
      feature-regex = "^user-.*$"
      contribution-regex = "^renovate/"
      observed-regex = "^dependabot/"

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
        perennial regex: ^release-
        default branch type: observed
        feature regex: ^user-.*$
        parked branches: (none)
        contribution branches: (none)
        contribution regex: ^renovate/
        observed branches: (none)
        observed regex: ^dependabot/

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: yes
        ship strategy: squash-merge
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
    Given the main branch is "git-main"
    And Git Town setting "perennial-branches" is "git-perennial-1 git-perennial-2"
    And Git Town setting "observed-branches" is "observed-1 observed-2"
    And Git Town setting "contribution-branches" is "contribution-1 contribution-2"
    And Git Town setting "parked-branches" is "parked-1 parked-2"
    And Git Town setting "perennial-regex" is "^git-perennial-"
    And Git Town setting "feature-regex" is "git-feature-.*"
    And Git Town setting "default-branch-type" is "observed"
    And Git Town setting "push-new-branches" is "false"
    And Git Town setting "ship-strategy" is "squash-merge"
    And Git Town setting "ship-delete-tracking-branch" is "false"
    And Git Town setting "sync-upstream" is "false"
    And Git Town setting "sync-tags" is "false"
    And Git Town setting "sync-perennial-strategy" is "merge"
    And Git Town setting "sync-feature-strategy" is "merge"
    And the configuration file:
      """
      push-new-branches = true
      ship-strategy = "api"
      ship-delete-tracking-branch = true
      sync-upstream = true
      sync-tags = true

      [branches]
      main = "config-main"
      perennials = [ "config-perennial-1", "config-perennial-2" ]
      perennial-regex = "^config-perennial-"
      default-type = "contribution"
      feature-regex = "^config-feature-.*$"

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
        perennial regex: ^git-perennial-
        default branch type: observed
        feature regex: git-feature-.*
        parked branches: parked-1, parked-2
        contribution branches: contribution-1, contribution-2
        observed branches: observed-1, observed-2

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship strategy: squash-merge
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
      | NAME   | TYPE      | PARENT | LOCATIONS |
      | alpha  | feature   | main   | local     |
      | qa     | perennial |        | local     |
      | beta   | feature   | main   | local     |
      | child  | feature   | alpha  | local     |
      | hotfix | feature   | qa     | local     |
    When I run "git-town config"
    Then it prints:
      """
      Branches:
        main branch: main
        perennial branches: qa
        perennial regex: (not set)
        default branch type: feature
        feature regex: (not set)
        parked branches: (none)
        contribution branches: (none)
        observed branches: (none)

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship strategy: api
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
        default branch type: feature
        feature regex: (not set)
        parked branches: (none)
        contribution branches: (none)
        observed branches: (none)

      Configuration:
        offline: no
        run pre-push hook: yes
        push new branches: no
        ship strategy: api
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
