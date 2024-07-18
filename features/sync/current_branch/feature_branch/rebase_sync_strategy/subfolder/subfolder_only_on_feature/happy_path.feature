Feature: sync inside a folder that doesn't exist on the main branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
    Given Git Town setting "sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME        |
      | main   | local, origin | main commit   | main_file        |
      | alpha  | local, origin | folder commit | new_folder/file1 |
      | beta   | local, origin | beta commit   | file2            |
    And the current branch is "alpha"
    And an uncommitted file
    When I run "git-town sync --all" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git fetch --prune --tags                        |
      |        | git add -A                                      |
      |        | git stash                                       |
      |        | git checkout main                               |
      | main   | git rebase origin/main                          |
      |        | git checkout alpha                              |
      | alpha  | git rebase main                                 |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git rebase main                                 |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout alpha                              |
      | alpha  | git push --tags                                 |
      |        | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "alpha"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | alpha  | local, origin | main commit   |
      |        |               | folder commit |
      | beta   | local, origin | main commit   |
      |        |               | beta commit   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                               |
      | alpha  | git add -A                                            |
      |        | git stash                                             |
      |        | git reset --hard {{ sha-before-run 'folder commit' }} |
      |        | git push --force-with-lease --force-if-includes       |
      |        | git checkout beta                                     |
      | beta   | git reset --hard {{ sha-before-run 'beta commit' }}   |
      |        | git push --force-with-lease --force-if-includes       |
      |        | git checkout alpha                                    |
      | alpha  | git stash pop                                         |
    And the current branch is still "alpha"
    And the initial commits exist
    And the initial branches and lineage exist
