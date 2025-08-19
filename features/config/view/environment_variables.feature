Feature: display configuration defined in environment variables

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

  @this
  Scenario: all configured in Git, no stacked changes
    When I run "git-town config" with these environment variables
      | GIT_TOWN_PERENNIAL_BRANCHES  | qa staging   |
      | GIT_TOWN_PERENNIAL_REGEX     | ^release-    |
      | GIT_TOWN_CONTRIBUTION_REGEX  | ^renovate/   |
      | GIT_TOWN_OBSERVED_REGEX      | ^dependabot/ |
      | GIT_TOWN_FEATURE_REGEX       | ^user-.*$    |
      | GIT_TOWN_SHIP_STRATEGY       | squash-merge |
      | GIT_TOWN_UNKNOWN_BRANCH_TYPE | observed     |
      | GIT_TOWN_OFFLINE             |            1 |
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
        offline: yes
      
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
        GitLab connector type: (not set)
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
