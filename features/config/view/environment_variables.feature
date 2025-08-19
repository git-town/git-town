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
      | GIT_TOWN_PERENNIAL_BRANCHES          | qa staging         |
      | GIT_TOWN_PERENNIAL_REGEX             | ^release-          |
      | GIT_TOWN_CONTRIBUTION_REGEX          | ^renovate/         |
      | GIT_TOWN_OBSERVED_REGEX              | ^dependabot/       |
      | GIT_TOWN_FEATURE_REGEX               | ^user-.*$          |
      | GIT_TOWN_SHIP_STRATEGY               | squash-merge       |
      | GIT_TOWN_UNKNOWN_BRANCH_TYPE         | observed           |
      | GIT_TOWN_OFFLINE                     |                  1 |
      | GIT_TOWN_NEW_BRANCH_TYPE             | prototype          |
      | GIT_TOWN_SHARE_NEW_BRANCHES          | push               |
      | GIT_TOWN_FORGE_TYPE                  | gitlab             |
      | GIT_TOWN_DEV_REMOTE                  | my-fork            |
      | GIT_TOWN_ORIGIN_HOSTNAME             | codeforge          |
      | GIT_TOWN_BITBUCKET_USERNAME          | bitbucket-user     |
      | GIT_TOWN_BITBUCKET_APP_PASSWORD      | bitbucket-password |
      | GIT_TOWN_CODEBERG_TOKEN              | codeberg-token     |
      | GIT_TOWN_GITEA_TOKEN                 | gitea-token        |
      | GIT_TOWN_GITHUB_CONNECTOR_TYPE       | gh                 |
      | GIT_TOWN_GITHUB_TOKEN                | github-token       |
      | GIT_TOWN_GITLAB_CONNECTOR_TYPE       | glab               |
      | GIT_TOWN_GITLAB_TOKEN                | gitlab-token       |
      | GIT_TOWN_SHIP_DELETE_TRACKING_BRANCH |                  0 |
      | GIT_TOWN_SHIP_STRATEGY               | fast-forward       |
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
        new branch type: prototype
        share new branches: push
      
      Hosting:
        development remote: my-fork
        forge type: gitlab
        origin hostname: codeforge
        Bitbucket username: bitbucket-user
        Bitbucket app password: bitbucket-password
        Codeberg token: codeberg-token
        Gitea token: gitea-token
        GitHub connector type: gh
        GitHub token: github-token
        GitLab connector type: glab
        GitLab token: gitlab-token
      
      Ship:
        delete tracking branch: no
        ship strategy: fast-forward
      
      Sync:
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: rebase
        prototype sync strategy: merge
        sync tags: yes
        sync with upstream: yes
      """
