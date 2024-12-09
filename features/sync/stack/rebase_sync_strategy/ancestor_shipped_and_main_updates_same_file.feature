Feature: shipped the head branch of a synced stack with dependent changes that create a file while main also creates the same file

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | file      | alpha content |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta commit | file      | beta content |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And origin ships the "alpha" branch using the "squash-merge" ship-strategy
    And I add this commit to the "main" branch
      | MESSAGE                    | FILE NAME | FILE CONTENT |
      | independent commit on main | file      | main content |
    And the current branch is "beta"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | beta   | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "file" with "resolved main content"
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | main   | git -c core.editor=true rebase --continue |
      |        | git push                                  |
      |        | git checkout beta                         |
      | beta   | git pull                                  |
      |        | git rebase --onto main alpha              |
      |        | git checkout --theirs file                |
      |        | git add file                              |
      |        | git -c core.editor=true rebase --continue |
      |        | git push --force-with-lease               |
      |        | git branch -D alpha                       |
    And all branches are now synchronized
    And the current branch is now "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                    | FILE NAME | FILE CONTENT          |
      | main   | local, origin | alpha commit               | file      | alpha content         |
      |        |               | independent commit on main | file      | resolved main content |
      | beta   | local, origin | beta commit                | file      | beta content          |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git checkout beta  |
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                    | FILE NAME | FILE CONTENT  |
      | main   | local         | independent commit on main | file      | main content  |
      |        | origin        | alpha commit               | file      | alpha content |
      | alpha  | local         | alpha commit               | file      | alpha content |
      | beta   | local, origin | beta commit                | file      | beta content  |
      |        | origin        | alpha commit               | file      | alpha content |
    And the initial branches and lineage exist now
