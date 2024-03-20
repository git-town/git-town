Feature: handle rebase conflicts between perennial branch and its tracking branch

  Background:
    Given the perennial branches "alpha", "beta", and "gamma"
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT        |
      | main   | origin        | main commit        | main_file        | main content        |
      | alpha  | local, origin | alpha commit       | alpha_file       | alpha content       |
      | beta   | local         | local beta commit  | conflicting_file | local beta content  |
      |        | origin        | origin beta commit | conflicting_file | origin beta content |
      | gamma  | local, origin | gamma commit       | gamma_file       | gamma content       |
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then I am not prompted for any parent branches
    And it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git checkout alpha       |
      | alpha  | git rebase origin/alpha  |
      |        | git checkout beta        |
      | beta   | git rebase origin/beta   |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      To continue by skipping the current branch, run "git-town skip".
      """
    And the uncommitted file is stashed
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND            |
      | beta   | git rebase --abort |
      |        | git checkout main  |
      | main   | git stash pop      |
    And the current branch is now "main"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: skip
    When I run "git-town skip"
    Then it runs the commands
      | BRANCH | COMMAND                 |
      | beta   | git rebase --abort      |
      |        | git checkout gamma      |
      | gamma  | git rebase origin/gamma |
      |        | git checkout main       |
      | main   | git rebase origin/main  |
      |        | git push --tags         |
      |        | git stash pop           |
    And the current branch is now "main"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            |
      | main   | local, origin | main commit        |
      | alpha  | local, origin | alpha commit       |
      | beta   | local         | local beta commit  |
      |        | origin        | origin beta commit |
      | gamma  | local, origin | gamma commit       |

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the uncommitted file is stashed
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND                 |
      | beta   | git rebase --continue   |
      |        | git push                |
      |        | git checkout gamma      |
      | gamma  | git rebase origin/gamma |
      |        | git checkout main       |
      | main   | git rebase origin/main  |
      |        | git push --tags         |
      |        | git stash pop           |
    And all branches are now synchronized
    And the current branch is now "main"
    And the uncommitted file still exists
    And no rebase is in progress

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH | COMMAND                 |
      | beta   | git push                |
      |        | git checkout gamma      |
      | gamma  | git rebase origin/gamma |
      |        | git checkout main       |
      | main   | git rebase origin/main  |
      |        | git push --tags         |
      |        | git stash pop           |
