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
      | feature regex               | enter             |                                             |
      | unknown branch type         | enter             |                                             |
      | dev-remote                  | enter             |                                             |
      | origin hostname             | enter             |                                             |
      | forge type: auto-detect     | enter             |                                             |
      | gitlab connector type: api  | enter             |                                             |
      | gitlab token                | 1 2 3 4 5 6 enter |                                             |
      | token scope                 | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-prototype-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | sync-tags                   | enter             |                                             |
      | share-new-branches          | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | new-branch-type             | enter             |                                             |
      | ship-strategy               | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config --local git-town.gitlab-token 123456 |
      | git config git-town.new-branch-type feature     |
      | git config git-town.gitlab-connector api        |
    And local Git setting "git-town.forge-type" still doesn't exist

  Scenario: select GitLab manually
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
      | forge type                  | up enter          |                                             |
      | gitlab connector type: api  | enter             |                                             |
      | gitlab token                | 1 2 3 4 5 6 enter |                                             |
      | token scope                 | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-prototype-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | sync-tags                   | enter             |                                             |
      | share-new-branches          | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | new-branch-type             | enter             |                                             |
      | ship-strategy               | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config --local git-town.gitlab-token 123456 |
      | git config git-town.new-branch-type feature     |
      | git config git-town.forge-type gitlab           |
      | git config git-town.gitlab-connector api        |
    And local Git setting "git-town.forge-type" is now "gitlab"
    And local Git setting "git-town.gitlab-token" is now "123456"

  Scenario: store GitLab API token globally
    Given my repo's "origin" remote is "git@gitlab.com:git-town/git-town.git"
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
      | forge type                  | enter             |                                             |
      | gitlab connector type: api  | enter             |                                             |
      | gitlab token                | 1 2 3 4 5 6 enter |                                             |
      | token scope                 | down enter        |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-prototype-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | sync-tags                   | enter             |                                             |
      | share-new-branches          | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | new-branch-type             | enter             |                                             |
      | ship-strategy               | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |
    Then Git Town runs the commands
      | COMMAND                                          |
      | git config --global git-town.gitlab-token 123456 |
      | git config git-town.new-branch-type feature      |
      | git config git-town.gitlab-connector api         |
    And global Git setting "git-town.gitlab-token" is now "123456"

  Scenario: edit global GitLab token
    Given my repo's "origin" remote is "git@gitlab.com:git-town/git-town.git"
    And global Git setting "git-town.gitlab-token" is "123"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS                                      | DESCRIPTION                                 |
      | welcome                     | enter                                     |                                             |
      | aliases                     | enter                                     |                                             |
      | main branch                 | enter                                     |                                             |
      | perennial branches          |                                           | no input here since the dialog doesn't show |
      | perennial regex             | enter                                     |                                             |
      | feature regex               | enter                                     |                                             |
      | unknown branch type         | enter                                     |                                             |
      | dev-remote                  | enter                                     |                                             |
      | origin hostname             | enter                                     |                                             |
      | forge type                  | enter                                     |                                             |
      | gitlab connector type: api  | enter                                     |                                             |
      | gitlab token                | backspace backspace backspace 4 5 6 enter |                                             |
      | token scope                 | enter                                     |                                             |
      | sync-feature-strategy       | enter                                     |                                             |
      | sync-perennial-strategy     | enter                                     |                                             |
      | sync-prototype-strategy     | enter                                     |                                             |
      | sync-upstream               | enter                                     |                                             |
      | sync-tags                   | enter                                     |                                             |
      | share-new-branches          | enter                                     |                                             |
      | push-hook                   | enter                                     |                                             |
      | new-branch-type             | enter                                     |                                             |
      | ship-strategy               | enter                                     |                                             |
      | ship-delete-tracking-branch | enter                                     |                                             |
      | save config to Git metadata | down enter                                |                                             |
    Then Git Town runs the commands
      | COMMAND                                       |
      | git config --global git-town.gitlab-token 456 |
      | git config git-town.new-branch-type feature   |
      | git config git-town.gitlab-connector api      |
    And global Git setting "git-town.gitlab-token" is now "456"
