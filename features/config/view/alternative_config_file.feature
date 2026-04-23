Feature: show the configuration when using an alternative config file

  Scenario: all configured in config file with alternative filename
    Given a Git repo with origin
    And an uncommitted file ".git-branches.toml" with content:
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
        feature regex: (not set)
        main branch: main
        observed branches: (none)
        observed regex: (not set)
        parked branches: (none)
        perennial branches: public, staging
        perennial regex: (not set)
        prototype branches: (none)
        unknown branch type: feature
        order: asc
        display types: all branch types except "feature" and "main"

      Configuration:
        offline: no
        git user name: user
        git user email: email@example.com

      Create:
        branch prefix: (not set)
        new branch type: (not set)
        share new branches: no
        stash uncommitted changes: yes

      Hosting:
        browser: (not set)
        development remote: origin
        forge type: (not set)
        origin hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Forgejo token: (not set)
        Gitea token: (not set)
        GitHub connector: (not set)
        GitHub token: (not set)
        GitLab connector: (not set)
        GitLab token: (not set)

      Propose:
        breadcrumb: none
        breadcrumb direction: down

      Ship:
        delete tracking branch: yes
        ignore uncommitted changes: no
        ship strategy: api

      Sync:
        auto-resolve phantom conflicts: yes
        auto-sync: yes
        run detached: no
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: rebase
        prototype sync strategy: merge
        push branches: yes
        sync tags: yes
        sync with upstream: yes
      """
