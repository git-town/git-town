Feature: show the configuration from a subfolder

  Scenario: running in a subfolder
    Given a Git repo with origin
    And the configuration file:
      """
      [branches]
      main = "main"
      perennials = ["public"]
      """
    And a folder "subfolder"
    When I run "git-town config" in the "subfolder" folder
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
        perennial branches: public
        perennial regex: (not set)
        prototype branches: (none)

      Configuration:
        offline: no

      Create:
        new branch type: (not set)
        push new branches: no

      Hosting:
        development remote: origin
        hosting platform: (not set)
        hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Gitea token: (not set)
        GitHub token: (not set)
        GitLab token: (not set)

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
