Feature: sync inside a folder that doesn't exist on the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME        |
      | main   | local, origin | main commit   | main_file        |
      | alpha  | local, origin | folder commit | new_folder/file1 |
      | beta   | local, origin | beta commit   | file2            |
    And the current branch is "alpha"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town sync --all" in the "new_folder" folder

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git fetch --prune --tags                        |
      |        | git rebase main --no-update-refs                |
      |        | git rebase origin/alpha --no-update-refs        |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git rebase main --no-update-refs                |
      |        | git rebase origin/beta --no-update-refs         |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout alpha                              |
      | alpha  | git push --tags                                 |
    And all branches are now synchronized
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |
      | beta   | main   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | alpha  | local, origin | folder commit |
      |        |               | main commit   |
      | beta   | local, origin | beta commit   |
      |        |               | main commit   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                               |
      | alpha  | git reset --hard {{ sha-before-run 'folder commit' }} |
      |        | git push --force-with-lease --force-if-includes       |
      |        | git checkout beta                                     |
      | beta   | git reset --hard {{ sha-before-run 'beta commit' }}   |
      |        | git push --force-with-lease --force-if-includes       |
      |        | git checkout alpha                                    |
    And the initial commits exist now
    And the initial branches and lineage exist now
