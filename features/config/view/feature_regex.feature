@smoke
Feature: with feature regex

  Background:
    Given a Git repo with origin

  Scenario: configured in Git metadata
    Given the branches
      | NAME         | TYPE         | PARENT | LOCATIONS |
      | user-one     | (none)       | main   | local     |
      | qa           | perennial    |        | local     |
      | observed     | observed     |        | local     |
      | contribution | contribution |        | local     |
      | parked       | parked       | main   | local     |
    And local Git Town setting "feature-regex" is "user-.*"
    And local Git Town setting "default-branch-type" is "observed"
    When I run "git-town config"
    Then it prints:
      """
      Branches:
        main branch: main
        perennial branches: qa
        perennial regex: (not set)
        default branch type: observed branch
        feature regex: user-.*
        parked branches: parked
        contribution branches: contribution
        observed branches: observed

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

  Scenario: configured in config file
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
      perennial-regex = "release-.*"
      feature-regex = "user-.*"
      default-type = "observed"

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
        default branch type: observed branch
        feature regex: user-.*
        parked branches: (none)
        contribution branches: (none)
        observed branches: (none)

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
