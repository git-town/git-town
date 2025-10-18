@messyoutput
Feature: beam multiple commits from an unknown branch

  Background:
    Given a Git repo with origin
    And I ran "git checkout -b branch-1"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | branch-1 | local    | branch commit 1 |
      | branch-1 | local    | branch commit 2 |
    And the current branch is "branch-1"
    When I run "git-town hack --beam branch-2"
    # When I run "git-town hack --beam branch-2" and enter into the dialog:
    #   | DIALOG          | KEYS                   |
    #   | parent branch   | enter                  |
    #   | commits to beam | space down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                       |
      | branch-1 | git checkout -b branch-2 main |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        branch-2
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                |
      | branch-2 | git checkout branch-1  |
      | branch-1 | git branch -D branch-2 |
    And the initial branches and lineage exist now
    And the initial commits exist now
