Feature: show the configuration

  Scenario: all configured, no nested branches
    Given the main branch is "main"
    And the perennial branches are "qa" and "staging"
    When I run "git-town config"
    Then it prints:
      """
      BRANCHES
      Main branch: main
      Perennial branches: qa, staging

      Pull branch strategy: xxx
      push using --no-verify: (not set)
      push new branches: xxx
      ship deletes the remote branch:
      sync strategy: [ ] merge [x] rebase
      sync upstream:

      Hosting service: xxx

      Github token: (not set)
      Gitlab token: (not set)
      Gitea token: (not set)
      Offline: disabled

      Repo URL:
      Upstream URL:

      Aliases: (not set)
      """


  Scenario: all configured, with nested branches
    Given the perennial branches "qa" and "staging"
    And the feature branches "alpha" and "beta"
    And a feature branch "child" as a child of "alpha"
    And a feature branch "hotfix" as a child of "qa"
    When I run "git-town config"
    Then it prints:
      """
      Main branch: main
      Perennial branches: qa, staging

      Branch Ancestry:
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
      Main branch: [none]
      Perennial branches: [none]
      """
