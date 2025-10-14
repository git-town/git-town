Feature: display configuration from Git metadata

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE         | PARENT | LOCATIONS     |
      | contribution-1 | contribution |        | local, origin |
      | contribution-2 | contribution |        | local, origin |
      | observed-1     | observed     |        | local, origin |
      | observed-2     | observed     |        | local, origin |
      | parked-1       | parked       | main   | local         |
      | parked-2       | parked       | main   | local         |
      | perennial-1    | perennial    |        | local         |
      | perennial-2    | perennial    |        | local         |
      | prototype-1    | prototype    | main   | local         |
      | prototype-2    | prototype    | main   | local         |

  Scenario: all configured in Git, no stacked changes
    Given Git setting "git-town.perennial-branches" is "qa staging"
    And Git setting "git-town.perennial-regex" is "^release-"
    And Git setting "git-town.auto-sync" is "false"
    And Git setting "git-town.contribution-regex" is "^renovate/"
    And Git setting "git-town.observed-regex" is "^dependabot/"
    And Git setting "git-town.feature-regex" is "^user-.*$"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And Git setting "git-town.unknown-branch-type" is "observed"
    And Git setting "git-town.auto-resolve" is "false"
    And Git setting "git-town.detached" is "true"
    And Git setting "git-town.stash" is "false"
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: ^renovate/
        feature regex: ^user-.*$
        main branch: main
        observed branches: observed-1, observed-2
        observed regex: ^dependabot/
        parked branches: parked-1, parked-2
        perennial branches: qa, staging
        perennial regex: ^release-
        prototype branches: prototype-1, prototype-2
        unknown branch type: observed

      Configuration:
        offline: no

      Create:
        new branch type: (not set)
        share new branches: no
        stash uncommitted changes: no

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

      Ship:
        delete tracking branch: yes
        ship strategy: squash-merge

      Sync:
        auto-resolve phantom conflicts: no
        auto-sync: no
        run detached: yes
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: rebase
        prototype sync strategy: merge
        push branches: yes
        sync tags: yes
        sync with upstream: yes
      """

  Scenario: all configured, with stacked changes
    Given the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | alpha  | feature | main   | local     |
      | qa     | (none)  |        | local     |
      | beta   | feature | main   | local     |
      | child  | feature | alpha  | local     |
      | hotfix | feature | qa     | local     |
    And Git setting "git-town.perennial-branches" is "qa"
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: (not set)
        feature regex: (not set)
        main branch: main
        observed branches: observed-1, observed-2
        observed regex: (not set)
        parked branches: parked-1, parked-2
        perennial branches: qa
        perennial regex: (not set)
        prototype branches: prototype-1, prototype-2
        unknown branch type: feature

      Configuration:
        offline: no

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
        auto-resolve phantom conflicts: yes

      Branch Lineage:
        main
          alpha
            child
          beta
          parked-1
          parked-2
          prototype-1
          prototype-2

        qa
          hotfix
      """
