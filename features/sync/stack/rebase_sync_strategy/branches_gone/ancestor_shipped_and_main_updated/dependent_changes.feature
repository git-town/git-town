Feature: shipped the head branch of a synced stack with dependent changes that create a file while main also creates the same file

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
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
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT           |
      | beta   | local, origin | beta commit | file      | alpha and beta content |
    And origin ships the "alpha" branch using the "squash-merge" ship-strategy
    And I add this commit to the "main" branch
      | MESSAGE                    | FILE NAME | FILE CONTENT |
      | independent commit on main | file      | main content |
    And the current branch is "beta"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | beta   | git fetch --prune --tags                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And a rebase is now in progress
    And file "file" now has content:
      """
      <<<<<<< HEAD
      alpha content
      =======
      main content
      >>>>>>> {{ sha-short 'independent commit on main' }} (independent commit on main)
      """

  Scenario: resolve and continue
    When I resolve the conflict in "file" with "alpha and independent content"
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                 |
      | main   | GIT_EDITOR=true git rebase --continue                   |
      |        | git push                                                |
      |        | git checkout beta                                       |
      | beta   | git pull                                                |
      |        | git -c rebase.updateRefs=false rebase --onto main alpha |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And a rebase is now in progress
    And file "file" now has content:
      """
      <<<<<<< HEAD
      alpha and independent content
      =======
      alpha and beta content
      >>>>>>> {{ sha-short 'beta commit' }} (beta commit)
      """
    When I resolve the conflict in "file" with "alpha, beta, and independent content"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | beta   | GIT_EDITOR=true git rebase --continue |
      |        | git push --force-with-lease           |
      |        | git branch -D alpha                   |
    And no rebase is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                    | FILE NAME | FILE CONTENT                         |
      | main   | local, origin | alpha commit               | file      | alpha content                        |
      |        |               | independent commit on main | file      | alpha and independent content        |
      | beta   | local, origin | beta commit                | file      | alpha, beta, and independent content |

  @this
  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git checkout beta  |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                    | FILE NAME | FILE CONTENT           |
      | main   | local         | independent commit on main | file      | main content           |
      |        | origin        | alpha commit               | file      | alpha content          |
      | alpha  | local         | alpha commit               | file      | alpha content          |
      | beta   | local, origin | beta commit                | file      | alpha and beta content |
      |        | origin        | alpha commit               | file      | alpha content          |
    And the initial branches and lineage exist now
