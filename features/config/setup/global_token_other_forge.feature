@messyoutput
Feature: a global API token of another forge exists

  Scenario: on GitHub, with global GitLab token
    Given a Git repo with origin
    And my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    And global Git setting "git-town.gitlab-token" is "987654"
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
      | github connector type       | enter           |                                             |
      | github token                | g h t o k enter |                                             |
      | token scope                 | enter           |                                             |
      | sync feature strategy       | enter           |                                             |
      | sync perennial strategy     | enter           |                                             |
      | sync prototype strategy     | enter           |                                             |
      | sync upstream               | enter           |                                             |
      | sync tags                   | enter           |                                             |
      | detached                    | enter           |                                             |
      | stash                       | enter           |                                             |
      | share new branches          | enter           |                                             |
      | push hook                   | enter           |                                             |
      | ship strategy               | enter           |                                             |
      | ship delete tracking branch | enter           |                                             |
      | config storage              | enter           |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.github-token ghtok               |
      | git config git-town.detached false                   |
      | git config git-town.new-branch-type feature          |
      | git config git-town.github-connector api             |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.stash true                       |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy ff-only  |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" is now "ghtok"
    And local Git setting "git-town.gitlab-token" now doesn't exist
    And global Git setting "git-town.gitlab-token" is still "987654"
