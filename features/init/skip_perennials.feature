@messyoutput
Feature: don't ask for perennial branches if no branches that could be perennial exist

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS       | DESCRIPTION                                 |
      | welcome            | enter      |                                             |
      | aliases            | enter      |                                             |
      | main branch        | down enter |                                             |
      | perennial branches |            | no input here since the dialog doesn't show |
      | origin hostname    | enter      |                                             |
      | forge type         | enter      |                                             |
      | enter all          | enter      |                                             |
      | config storage     | enter      |                                             |

  Scenario: result
    Then the main branch is now "main"
    And there are still no perennial branches
