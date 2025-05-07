Feature: show the configuration in Spanish

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE         | PARENT | LOCATIONS |
      | contribution-1 | contribution |        | local     |
      | contribution-2 | contribution |        | local     |
      | observed-1     | observed     |        | local     |
      | observed-2     | observed     |        | local     |
      | parked-1       | parked       | main   | local     |
      | parked-2       | parked       | main   | local     |
      | perennial-1    | perennial    |        | local     |
      | perennial-2    | perennial    |        | local     |
      | prototype-1    | prototype    | main   | local     |
      | prototype-2    | prototype    | main   | local     |
    And Git setting "git-town.perennial-branches" is "qa staging"
    And Git setting "git-town.perennial-regex" is "^release-"
    And Git setting "git-town.contribution-regex" is "^renovate/"
    And Git setting "git-town.observed-regex" is "^dependabot/"
    And Git setting "git-town.default-branch-type" is "observed"
    And Git setting "git-town.feature-regex" is "^user-.*$"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    When I run "git-town config" with these environment variables
      | LANG | es_ES.UTF-8 |

  Scenario: result
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: ^renovate/
        default branch type: observed
        feature regex: ^user-.*$
        main branch: main
        observed branches: observed-1, observed-2
        observed regex: ^dependabot/
        parked branches: parked-1, parked-2
        perennial branches: qa, staging
        perennial regex: ^release-
        prototype branches: prototype-1, prototype-2

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
        GitHub token: (not set)
        GitLab token: (not set)

      Ship:
        delete tracking branch: yes
        ship strategy: squash-merge

      Sync:
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: rebase
        prototype sync strategy: merge
        sync tags: yes
        sync with upstream: yes
      """
