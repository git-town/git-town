Feature: handle rebase conflicts between perennial branch and its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE      | LOCATIONS     |
      | alpha | perennial | local, origin |
      | beta  | perennial | local, origin |
      | gamma | perennial | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT        |
      | main   | origin        | main commit        | main_file        | main content        |
      | alpha  | local, origin | alpha commit       | alpha_file       | alpha content       |
      | beta   | local         | local beta commit  | conflicting_file | local beta content  |
      |        | origin        | origin beta commit | conflicting_file | origin beta content |
      | gamma  | local, origin | gamma commit       | gamma_file       | gamma content       |
    And the current branch is "main"
    When I run "git-town sync --all"

  Scenario: result
    Then I am not prompted for any parent branches
    And Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git fetch --prune --tags                          |
      |        | git checkout beta                                 |
      | beta   | git -c rebase.updateRefs=false rebase origin/beta |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND            |
      | beta   | git rebase --abort |
      |        | git checkout main  |
    And the initial branches and lineage exist now
    And the initial commits exist now

  Scenario: skip
    When I run "git-town skip"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | beta   | git rebase --abort                                |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git push --tags                                   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            |
      | main   | local, origin | main commit        |
      | alpha  | local, origin | alpha commit       |
      | beta   | local         | local beta commit  |
      |        | origin        | origin beta commit |
      | gamma  | local, origin | gamma commit       |

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | beta   | GIT_EDITOR=true git rebase --continue             |
      |        | git push                                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git push --tags                                   |
    And no rebase is now in progress
    And all branches are now synchronized

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | beta   | git push                                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git push --tags                                   |
