Feature: override an existing Git alias

  Background:
    Given I ran "git config --global alias.append checkout"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS    |
      | welcome                     | enter   |
      | aliases                     | o enter |
      | main development branch     | enter   |
      | perennial branches          | enter   |
      | perennial regex             | enter   |
      | hosting platform            | enter   |
      | origin hostname             | enter   |
      | sync-feature-strategy       | enter   |
      | sync-perennial-strategy     | enter   |
      | sync-upstream               | enter   |
      | push-new-branches           | enter   |
      | push-hook                   | enter   |
      | ship-delete-tracking-branch | enter   |
      | sync-before-ship            | enter   |
      | save config to config file  | enter   |

  Scenario: result
    Then it runs the commands
      | COMMAND                                        |
      | git config --global alias.append "town append" |
    And global Git setting "alias.append" is now "town append"
