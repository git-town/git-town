Feature: show the configuration

  Scenario: all configured, no nested branches
    Given the main branch is "main"
    And the perennial branches are "qa" and "staging"
    When I run "git-town config"
    Then it prints:
      """
      Branches:
        main branch: main
        perennial branches: qa, staging

      Configuration:
        offline: no
        pull branch strategy: rebase
        push using --no-verify: no
        push new branches: no
        ship removes the remote branch: yes
        sync strategy: merge
        sync with upstream: yes

      Hosting:
        hosting service: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """

  Scenario: all configured, with nested branches
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

      Configuration:
        offline: no
        pull branch strategy: rebase
        push using --no-verify: no
        push new branches: no
        ship removes the remote branch: yes
        sync strategy: merge
        sync with upstream: yes

      Hosting:
        hosting service: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)

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
      Branches:
        main branch: (not set)
        perennial branches: (not set)

      Configuration:
        offline: no
        pull branch strategy: rebase
        push using --no-verify: no
        push new branches: no
        ship removes the remote branch: yes
        sync strategy: merge
        sync with upstream: yes

      Hosting:
        hosting service: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """
