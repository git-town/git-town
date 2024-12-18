Feature: show the configuration when using an alternative config file

  Scenario: all configured in config file with alternative filename
    Given a Git repo with origin
    And file ".git-town.toml" with content
      """
      [branches]
      main = "main"
      perennials = [ "public", "staging" ]
      """
    When I run "git-town config"
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
        perennial branches: public, staging
        perennial regex: (not set)

      Configuration:
        offline: no

      Create:
        new branch type: feature
        push new branches: no

      Hosting:
        hosting platform: (not set)
        hostname: (not set)
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)

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
