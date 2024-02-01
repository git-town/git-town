Feature: aborting the setup assistant

  Background:
    And local Git setting "init.defaultbranch" is "main"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                  | KEYS  |
      | welcome                 | enter |
      | aliases                 | enter |
      | main development branch | enter |
      | perennial branches      | enter |
      | hosting platform        | esc   |

  Scenario: result
    Then it runs no commands
    And the main branch is still not set
    And there are still no perennial branches
    And local Git Town setting "code-hosting-platform" is still not set
    And local Git Town setting "push-new-branches" is still not set
    And local Git Town setting "push-hook" is still not set
    And local Git Town setting "sync-feature-strategy" is still not set
    And local Git Town setting "sync-perennial-strategy" is still not set
    And local Git Town setting "sync-upstream" is still not set
    And local Git Town setting "ship-delete-tracking-branch" is still not set
    And local Git Town setting "sync-before-ship" is still not set
