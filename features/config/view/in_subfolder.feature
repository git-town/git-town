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
        feature regex: (not set)
        main branch: main
        observed branches: (none)
        observed regex: (not set)
        parked branches: (none)
        perennial branches: public
        perennial regex: (not set)
        prototype branches: (none)
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
        GitHub connector type: (not set)
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
