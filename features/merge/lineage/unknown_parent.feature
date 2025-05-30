@messyoutput
Feature: merging with missing lineage

  Background:
    Given a Git repo with origin
    And I ran "git checkout -b alpha"
    And I ran "git checkout -b beta"
    When I run "git-town merge" and enter into the dialog:
      | DIALOG                          | KEYS       |
      | select parent branch for "beta" | down enter |
      | select parent branch for "beta" | enter      |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
      |        | git checkout alpha       |
      | alpha  | git branch -D beta       |
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                    |
      | alpha  | git branch beta {{ sha 'initial commit' }} |
      |        | git checkout beta                          |
    And the initial commits exist now
    And the initial lineage exists now
