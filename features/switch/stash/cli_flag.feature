@messyoutput
Feature: disable stashing via CLI flag

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | local-1 | feature | main   | local, origin |
      | local-2 | feature | main   | local, origin |
    And the current branch is "local-1"
    And an uncommitted file
    When I run "git-town switch --stash" and enter into the dialogs:
      | DIALOG        | KEYS       |
      | switch-branch | down enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                     |
      | local-1 | git add -A                  |
      |         | git stash -m "Git Town WIP" |
      |         | git checkout local-2        |
      | local-2 | git stash pop               |
