Feature: remove an existing code hosting override

  Background:
    Given local Git Town setting "code-hosting-platform" is "github"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS           | DESCRIPTION                                 |
      | welcome                     | enter          |                                             |
      | aliases                     | enter          |                                             |
      | main development branch     | down enter     |                                             |
      | perennial branches          |                | no input here since the dialog doesn't show |
      | hosting platform            | up up up enter |                                             |
      | origin hostname             | enter          |                                             |
      | sync-feature-strategy       | enter          |                                             |
      | sync-perennial-strategy     | enter          |                                             |
      | sync-upstream               | enter          |                                             |
      | push-new-branches           | enter          |                                             |
      | push-hook                   | enter          |                                             |
      | ship-delete-tracking-branch | enter          |                                             |
      | sync-before-ship            | enter          |                                             |
      | save config to Git metadata | down enter     |                                             |

  Scenario: result
    Then it runs the commands
      | COMMAND                                           |
      | git config --unset git-town.code-hosting-platform |
    And local Git Town setting "code-hosting-platform" is now not set
