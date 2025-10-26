Feature: no configuration data

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE         | PARENT | LOCATIONS     |
      | contribution-1 | contribution |        | local, origin |
      | contribution-2 | contribution |        | local, origin |
      | observed-1     | observed     |        | local, origin |
      | observed-2     | observed     |        | local, origin |
    And Git Town is not configured

  Scenario: no configuration data
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: (not set)
        feature regex: (not set)
        main branch: (not set)
        observed branches: observed-1, observed-2
        observed regex: (not set)
        parked branches: (none)
        perennial branches: (none)
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
        new branch type: (not set)
        share new branches: no
        stash uncommitted changes: yes

      Hosting:
        development remote: origin
        forge type: (not set)
        origin hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Forgejo token: (not set)
        Gitea token: (not set)
        GitHub connector type: (not set)
        GitHub token: (not set)
        GitLab connector type: (not set)
        GitLab token: (not set)

      Proposals:
        show lineage: none

      Ship:
        delete tracking branch: yes
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
