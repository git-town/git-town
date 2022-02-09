Feature: conflicts between the main branch and its tracking branch

  Background:
    Given a feature branch "existing"
    And the commits
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
      |        | origin   | conflicting origin commit | conflicting_file | origin content |
    And the current branch is "existing"
    And my workspace has an uncommitted file
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git add -A               |
      |          | git stash                |
      |          | git checkout main        |
      | main     | git rebase origin/main   |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And my repo now has a rebase in progress
    And my uncommitted file is stashed

  Scenario: abort
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | main     | git rebase --abort    |
      |          | git checkout existing |
      | existing | git stash pop         |
    And the current branch is now "existing"
    And my workspace has the uncommitted file again
    And there is no rebase in progress anymore
    And now the initial commits exist

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And my uncommitted file is stashed
    And my repo still has a rebase in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND               |
      | main   | git rebase --continue |
      |        | git push              |
      |        | git branch new main   |
      |        | git checkout new      |
      | new    | git stash pop         |
    And the current branch is now "new"
    And my workspace still contains my uncommitted file
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE                   |
      | main   | local, origin | conflicting origin commit |
      |        |               | conflicting local commit  |
      | new    | local         | conflicting origin commit |
      |        |               | conflicting local commit  |
    And my repo now has these committed files
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | resolved content |
      | new    | conflicting_file | resolved content |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH | COMMAND             |
      | main   | git push            |
      |        | git branch new main |
      |        | git checkout new    |
      | new    | git stash pop       |
    And the current branch is now "new"
    And my workspace still contains my uncommitted file
