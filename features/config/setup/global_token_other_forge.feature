@messyoutput
Feature: a global API token of another forge exists

  Scenario: on GitHub, with global GitLab token
    Given a Git repo with origin
    And my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    And global Git setting "git-town.gitlab-token" is "987654"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main branch                 | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | perennial regex             | enter             |                                             |
      | feature regex               | enter             |                                             |
      | unknown branch type         | enter             |                                             |
      | dev-remote                  | enter             |                                             |
      | origin hostname             | enter             |                                             |
      | forge type: auto-detect     | enter             |                                             |
      | github connector type: API  | enter             |                                             |
      | github token                | 1 2 3 4 5 6 enter |                                             |
      | token scope                 | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-prototype-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | sync-tags                   | enter             |                                             |
      | share-new-branches          | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | new-branch-type             | down enter        |                                             |
      | ship-strategy               | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config --local git-town.github-token 123456 |
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" is now "123456"
    And local Git setting "git-town.gitlab-token" now doesn't exist
    And global Git setting "git-town.gitlab-token" is still "987654"
