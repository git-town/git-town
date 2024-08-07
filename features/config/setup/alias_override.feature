@messyoutput
Feature: override an existing Git alias

  Background:
    Given a Git repo with origin
    And I ran "git config --global alias.append checkout"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS    |
      | welcome                     | enter   |
      | aliases                     | o enter |
      | main branch                 | enter   |
      | perennial branches          | enter   |
      | perennial regex             | enter   |
      | hosting platform            | enter   |
      | origin hostname             | enter   |
      | sync-feature-strategy       | enter   |
      | sync-perennial-strategy     | enter   |
      | sync-upstream               | enter   |
      | sync-tags                   | enter   |
      | push-new-branches           | enter   |
      | push-hook                   | enter   |
      | create-prototype-branches   | enter   |
      | ship-delete-tracking-branch | enter   |
      | save config to config file  | enter   |

  Scenario: result
    Then it runs the commands
      | COMMAND                                        |
      | git config --global alias.append "town append" |
    And global Git setting "alias.append" is now "town append"

  Scenario: undo
    When I run "git-town undo"
    Then global Git setting "alias.append" is now "checkout"
