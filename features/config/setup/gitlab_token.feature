@messyoutput
Feature: enter the GitLab API token

  Background:
    Given a Git repo with origin

  Scenario: auto-detected GitLab platform
    Given my repo's "origin" remote is "git@gitlab.com:git-town/git-town.git"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main branch                 | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | perennial regex             | enter             |                                             |
      | default branch type         | enter             |                                             |
      | feature regex               | enter             |                                             |
      | dev-remote                  | enter             |                                             |
      | forge type: auto-detect     | enter             |                                             |
      | gitlab token                | 1 2 3 4 5 6 enter |                                             |
      | origin hostname             | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-prototype-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | sync-tags                   | enter             |                                             |
      | push-new-branches           | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | new-branch-type             | enter             |                                             |
      | ship-strategy               | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |
    Then Git Town runs the commands
      | COMMAND                                 |
      | git config git-town.gitlab-token 123456 |
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.gitlab-token" is now "123456"

  Scenario: select GitLab manually
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main branch                 | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | perennial regex             | enter             |                                             |
      | default branch type         | enter             |                                             |
      | feature regex               | enter             |                                             |
      | dev-remote                  | enter             |                                             |
      | forge type                  | up enter          |                                             |
      | gitlab token                | 1 2 3 4 5 6 enter |                                             |
      | origin hostname             | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-prototype-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | sync-tags                   | enter             |                                             |
      | push-new-branches           | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | new-branch-type             | enter             |                                             |
      | ship-strategy               | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |
    Then Git Town runs the commands
      | COMMAND                                 |
      | git config git-town.gitlab-token 123456 |
      | git config git-town.forge-type gitlab   |
    And local Git setting "git-town.forge-type" is now "gitlab"
    And local Git setting "git-town.gitlab-token" is now "123456"

  Scenario: undo
    When I run "git-town undo"
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.gitlab-token" now doesn't exist
