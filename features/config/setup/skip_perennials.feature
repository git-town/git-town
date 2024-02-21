Feature: don't ask for perennial branches if no branches that could be perennial exist

  Background:
    Given Git Town is not configured
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS       | DESCRIPTION                                 |
      | welcome                     | enter      |                                             |
      | aliases                     | enter      |                                             |
      | main development branch     | down enter |                                             |
      | perennial branches          |            | no input here since the dialog doesn't show |
      | perennial regex             | enter      |                                             |
      | hosting platform            | enter      |                                             |
      | origin hostname             | enter      |                                             |
      | sync-feature-strategy       | enter      |                                             |
      | sync-perennial-strategy     | enter      |                                             |
      | sync-upstream               | enter      |                                             |
      | push-new-branches           | enter      |                                             |
      | push-hook                   | enter      |                                             |
      | ship-delete-tracking-branch | enter      |                                             |
      | sync-before-ship            | enter      |                                             |
      | save config to Git metadata | down enter |                                             |

  Scenario: result
    Then the main branch is now "main"
    And there are still no perennial branches
