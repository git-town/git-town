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
        main branch: main
        perennial branches: public
        perennial regex: (not set)
        default branch type: feature
        feature regex: (not set)
        parked branches: (none)
        contribution branches: (none)
        contribution regex: (none)
        observed branches: (none)
        observed regex: (none)

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
