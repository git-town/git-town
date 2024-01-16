Feature: enter Git Town configuration

  @this
  Scenario: change existing configuration
    Given a perennial branch "qa"
    And a branch "production"
    And the main branch is "main"
    And local Git Town setting "push-hook" is "false"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                   | KEYS                   | DESCRIPTION                               |
      | enter main branch        | enter                  | accept the already configured main branch |
      | enter perennial branches | space down space enter | configure the perennial branches          |
      | enter push-hook          | down enter             | disable the push hook                     |
    Then the main branch is now "main"
    And the perennial branches are now "production"
    And local Git Town setting "push-hook" is now "true"

  Scenario: unconfigured
    Given the branches "dev" and "production"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                  | KEYS                   |
      | main development branch | down enter             |
      | perennial branches      | space down space enter |
    Then the main branch is now "main"
    And the perennial branches are now "dev" and "production"

  Scenario: don't ask for perennial branches if no branches that could be perennial exist
    Given Git Town is not configured
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                  | KEYS       |
      | main development branch | down enter |
    Then the main branch is now "main"
    And there are still no perennial branches
