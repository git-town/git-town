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
      | contribution regex          | enter             |                                             |
      | observed regex              | enter             |                                             |
      | new branch type             | enter             |                                             |
      | unknown branch type         | enter             |                                             |
      | origin hostname             | enter             |                                             |
      | forge type                  | enter             |                                             |
      | gitlab connector type       | enter             |                                             |
      | gitlab token                | g l - t o k enter |                                             |
      | token scope                 | enter             |                                             |
      | sync feature strategy       | enter             |                                             |
      | sync perennial strategy     | enter             |                                             |
      | sync prototype strategy     | enter             |                                             |
      | sync upstream               | enter             |                                             |
      | sync tags                   | enter             |                                             |
      | share new branches          | enter             |                                             |
      | push hook                   | enter             |                                             |
      | ship strategy               | enter             |                                             |
      | ship delete tracking branch | enter             |                                             |
      | config storage              | enter             |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.gitlab-token gl-tok              |
      | git config git-town.new-branch-type feature          |
      | git config git-town.gitlab-connector api             |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy rebase   |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
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
      | contribution regex          | enter             |                                             |
      | observed regex              | enter             |                                             |
      | new branch type             | enter             |                                             |
      | unknown branch type         | enter             |                                             |
      | origin hostname             | enter             |                                             |
      | forge type                  | up enter          |                                             |
      | gitlab connector type       | enter             |                                             |
      | gitlab token                | g l - t o k enter |                                             |
      | token scope                 | enter             |                                             |
      | sync feature strategy       | enter             |                                             |
      | sync perennial strategy     | enter             |                                             |
      | sync prototype strategy     | enter             |                                             |
      | sync upstream               | enter             |                                             |
      | sync tags                   | enter             |                                             |
      | share new branches          | enter             |                                             |
      | push hook                   | enter             |                                             |
      | ship strategy               | enter             |                                             |
      | ship delete tracking branch | enter             |                                             |
      | config storage              | enter             |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.gitlab-token gl-tok              |
      | git config git-town.new-branch-type feature          |
      | git config git-town.forge-type gitlab                |
      | git config git-town.gitlab-connector api             |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy rebase   |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And local Git setting "git-town.forge-type" is now "gitlab"
    And local Git setting "git-town.gitlab-token" is now "gl-tok"

  Scenario: store GitLab API token globally
    Given my repo's "origin" remote is "git@gitlab.com:git-town/git-town.git"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS            | DESCRIPTION                                 |
      | welcome                     | enter           |                                             |
      | aliases                     | enter           |                                             |
      | main branch                 | enter           |                                             |
      | perennial branches          |                 | no input here since the dialog doesn't show |
      | perennial regex             | enter           |                                             |
      | feature regex               | enter           |                                             |
      | contribution regex          | enter           |                                             |
      | observed regex              | enter           |                                             |
      | new branch type             | enter           |                                             |
      | unknown branch type         | enter           |                                             |
      | origin hostname             | enter           |                                             |
      | forge type                  | enter           |                                             |
      | gitlab connector type       | enter           | api                                         |
      | gitlab token                | g l t o k enter |                                             |
      | token scope                 | down enter      |                                             |
      | sync feature strategy       | enter           |                                             |
      | sync perennial strategy     | enter           |                                             |
      | sync prototype strategy     | enter           |                                             |
      | sync upstream               | enter           |                                             |
      | sync tags                   | enter           |                                             |
      | share new branches          | enter           |                                             |
      | push hook                   | enter           |                                             |
      | ship strategy               | enter           |                                             |
      | ship delete tracking branch | enter           |                                             |
      | config storage              | enter           | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config --global git-town.gitlab-token gltok      |
      | git config git-town.new-branch-type feature          |
      | git config git-town.gitlab-connector api             |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy rebase   |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And global Git setting "git-town.gitlab-token" is now "gltok"

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
      | contribution regex          | enter                                     |                                             |
      | observed regex              | enter                                     |                                             |
      | new branch type             | enter                                     |                                             |
      | unknown branch type         | enter                                     |                                             |
      | origin hostname             | enter                                     |                                             |
      | forge type                  | enter                                     |                                             |
      | gitlab connector type       | enter                                     |                                             |
      | gitlab token                | backspace backspace backspace 4 5 6 enter |                                             |
      | token scope                 | enter                                     |                                             |
      | sync feature strategy       | enter                                     |                                             |
      | sync perennial strategy     | enter                                     |                                             |
      | sync prototype strategy     | enter                                     |                                             |
      | sync upstream               | enter                                     |                                             |
      | sync tags                   | enter                                     |                                             |
      | share new branches          | enter                                     |                                             |
      | push hook                   | enter                                     |                                             |
      | ship strategy               | enter                                     |                                             |
      | ship delete tracking branch | enter                                     |                                             |
      | config storage              | enter                                     |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config --global git-town.gitlab-token 456        |
      | git config git-town.new-branch-type feature          |
      | git config git-town.gitlab-connector api             |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy rebase   |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And global Git setting "git-town.gitlab-token" is now "456"
