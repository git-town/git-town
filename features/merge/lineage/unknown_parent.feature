@messyoutput
Feature: merging with missing lineage

  Background:
    Given a Git repo with origin
    And I ran "git checkout -b alpha"
    And I ran "git checkout -b beta"
    When I run "git-town merge" and enter into the dialog:
      | DIALOG                    | KEYS       |
      | parent branch for "beta"  | down enter |
      | parent branch for "alpha" | enter      |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
      |        | git checkout alpha       |
      | alpha  | git branch -D beta       |
    And this lineage exists now
      """
      main
        alpha
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                    |
      | alpha  | git branch beta {{ sha 'initial commit' }} |
      |        | git checkout beta                          |
    And the initial lineage exists now
    And the initial commits exist now
