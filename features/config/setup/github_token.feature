Feature: enter the GitHub API token

  Scenario: select GitHub manually
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS                 | DESCRIPTION                                 |
      | welcome                     | enter                |                                             |
      | aliases                     | enter                |                                             |
      | main development branch     | enter                |                                             |
      | perennial branches          |                      | no input here since the dialog doesn't show |
      | hosting platform            | down down down enter |                                             |
      | github token                | 1 2 3 4 5 6 enter    |                                             |
      | origin hostname             | enter                |                                             |
      | sync-feature-strategy       | enter                |                                             |
      | sync-perennial-strategy     | enter                |                                             |
      | sync-upstream               | enter                |                                             |
      | push-new-branches           | enter                |                                             |
      | push-hook                   | enter                |                                             |
      | ship-delete-tracking-branch | enter                |                                             |
      | sync-before-ship            | enter                |                                             |
      | save config to Git metadata | down enter           |                                             |
    Then it runs the commands
      | COMMAND                                          |
      | git config git-town.github-token 123456          |
      | git config git-town.code-hosting-platform github |
    And local Git Town setting "code-hosting-platform" is now "github"
    And local Git Town setting "github-token" is now "123456"
