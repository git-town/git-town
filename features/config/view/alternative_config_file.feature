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
        contribution regex: (not set)
        default branch type: feature
        feature regex: (not set)
        main branch: main
        observed branches: (none)
        observed regex: (not set)
        parked branches: (none)
        perennial branches: public, staging
        perennial regex: (not set)
        prototype branches: (none)

      Configuration:
        offline: no

      Create:
        new branch type: (not set)
        push new branches: no

      Hosting:
        development remote: origin
        forge type: (not set)
        origin hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Codeberg token: (not set)
        Gitea token: (not set)
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
