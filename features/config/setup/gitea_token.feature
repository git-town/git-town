Feature: enter the Gitea API token

  Scenario: auto-detected Gitea platform
    Given my repo's "origin" remote is "git@gitea.com:git-town/git-town.git"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                        | KEYS              | DESCRIPTION                                 |
      | welcome                       | enter             |                                             |
      | aliases                       | enter             |                                             |
      | main development branch       | enter             |                                             |
      | perennial branches            |                   | no input here since the dialog doesn't show |
      | hosting platform: auto-detect | enter             |                                             |
      | gitea token                   | 1 2 3 4 5 6 enter |                                             |
      | origin hostname               | enter             |                                             |
      | sync-feature-strategy         | enter             |                                             |
      | sync-perennial-strategy       | enter             |                                             |
      | sync-upstream                 | enter             |                                             |
      | push-new-branches             | enter             |                                             |
      | push-hook                     | enter             |                                             |
      | ship-delete-tracking-branch   | enter             |                                             |
      | sync-before-ship              | enter             |                                             |
      | save config to Git metadata   | down enter        |                                             |
    Then it runs the commands
      | COMMAND                                |
      | git config git-town.gitea-token 123456 |
    And local Git Town setting "code-hosting-platform" still doesn't exist
    And local Git Town setting "gitea-token" is now "123456"

  Scenario: select Gitea manually
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main development branch     | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | hosting platform            | down down enter   |                                             |
      | gitea token                 | 1 2 3 4 5 6 enter |                                             |
      | origin hostname             | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | push-new-branches           | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | sync-before-ship            | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |
    Then it runs the commands
      | COMMAND                                         |
      | git config git-town.gitea-token 123456          |
      | git config git-town.code-hosting-platform gitea |
    And local Git Town setting "code-hosting-platform" is now "gitea"
    And local Git Town setting "gitea-token" is now "123456"
