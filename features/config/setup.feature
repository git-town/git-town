Feature: enter Git Town configuration

  Scenario: unconfigured, accept all default values --> working setup
    Given the branches "dev" and "production"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                  | KEYS  |
      | main development branch | enter |
      | perennial branches      | enter |
      | enter push-new-branches | enter |
      | enter push-hook         | enter |
    Then the main branch is now "dev"
    And there are still no perennial branches
    And local Git Town setting "push-new-branches" is now "false"
    And local Git Town setting "push-hook" is now "true"

  Scenario: change existing configuration
    Given a perennial branch "qa"
    And a branch "production"
    And the main branch is "main"
    And local Git Town setting "push-new-branches" is "false"
    And local Git Town setting "push-hook" is "false"
    When I run "git-town config setup" and enter into the dialogs:
      | DESCRIPTION                               | KEYS                   |
      | accept the already configured main branch | enter                  |
      | configure the perennial branches          | space down space enter |
      | enable push-new-branches                  | down enter             |
      | disable the push hook                     | down enter             |
    Then the main branch is now "main"
    And the perennial branches are now "production"
    And local Git Town setting "push-new-branches" is now "true"
    And local Git Town setting "push-hook" is now "true"

  Scenario: don't ask for perennial branches if no branches that could be perennial exist
    Given Git Town is not configured
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                  | KEYS       | DESCRIPTION                                 |
      | main development branch | down enter |                                             |
      | perennial branches      |            | no input here since the dialog doesn't show |
      | enter push-new-branches | enter      |                                             |
      | enter push-hook         | enter      |                                             |
    Then the main branch is now "main"
    And there are still no perennial branches
