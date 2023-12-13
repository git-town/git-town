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
        run pre-push hook: yes
        push new branches: no
        ship deletes the tracking branch: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync with upstream: yes
        sync before shipping: no

      Hosting:
        hosting service override: (not set)
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
        run pre-push hook: yes
        push new branches: no
        ship deletes the tracking branch: yes
        sync-feature strategy: merge
        sync-perennial strategy: rebase
        sync with upstream: yes
        sync before shipping: no

      Hosting:
        hosting service override: (not set)
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
        perennial branches: (not set)

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
        hosting service override: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)
      """
