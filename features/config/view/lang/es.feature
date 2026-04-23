Feature: show the configuration in Spanish

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
    And Git setting "git-town.contribution-regex" is "^renovate/"
    And Git setting "git-town.feature-regex" is "^user-.*$"
    And Git setting "git-town.observed-regex" is "^dependabot/"
    And Git setting "git-town.perennial-branches" is "qa staging"
    And Git setting "git-town.perennial-regex" is "^release-"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And Git setting "git-town.unknown-branch-type" is "observed"
    When I run "git-town config" with these environment variables
      | LANG | es_ES.UTF-8 |

  Scenario: result
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
        ship strategy: squash-merge

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
