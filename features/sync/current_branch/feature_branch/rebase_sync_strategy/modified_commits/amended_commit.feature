Feature: rebase a branch that contains amended commits

# This test demonstrates how Git Town currently does not sync branches that contain amended commits correctly.
# This will be fixed by implementing https://github.com/git-town/git-town/issues/4586.

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
      | feature-1 | local, origin | commit 1a | file_1    | one          |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-2 | feature | feature-1 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature-2 | local, origin | commit 2 | file_2    | two          |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "feature-2"
    And wait 1 second to ensure new Git timestamps
    And I amend this commit
      | BRANCH    | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT |
      | feature-1 | local    | commit 1b | file_1    | another one  |
    And the current branch is "feature-2"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                      |
      | feature-2 | git fetch --prune --tags                     |
      |           | git checkout feature-1                       |
      | feature-1 | git rebase main --no-update-refs             |
      |           | git rebase origin/feature-1 --no-update-refs |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file_1
      """
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "file_1" with "another one"
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                         |
      | feature-1 | git -c core.editor=true rebase --continue       |
      |           | git push --force-with-lease --force-if-includes |
      |           | git checkout feature-2                          |
      | feature-2 | git rebase feature-1 --no-update-refs           |
      |           | git rebase origin/feature-2 --no-update-refs    |
      |           | git push --force-with-lease --force-if-includes |
    And the current branch is still "feature-2"
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
      | feature-1 | local, origin | commit 1a | file_1    | one          |
      |           |               | commit 1b | file_1    | another one  |
      | feature-2 | local, origin | commit 2  | file_2    | two          |
      |           |               | commit 1b | file_1    | another one  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                |
      | feature-1 | git rebase --abort     |
      |           | git checkout feature-2 |
    And the current branch is still "feature-2"
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE   |
      | feature-1 | local         | commit 1b |
      |           | origin        | commit 1a |
      | feature-2 | local         | commit 1a |
      |           | local, origin | commit 2  |
    And these branches exist now
      | REPOSITORY    | BRANCHES                   |
      | local, origin | main, feature-1, feature-2 |
