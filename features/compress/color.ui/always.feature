Feature: compress keeps the full commit message of the first commit

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE                              | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1\n\nbody line 1\nbody line 2 | file_1    | content 1    |
      |         |               | commit 2                             | file_2    | content 2    |
    And Git setting "color.ui" is "always"
    And the current branch is "feature"
    When I run "git-town compress"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git reset --soft main --                        |
      |         | git commit -m "commit 1                         |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE  |
      | feature | local, origin | commit 1 |
    And commit "commit 1" on branch "feature" now has this full commit message
      """
      commit 1

      body line 1
      body line 2
      """
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git reset --hard {{ sha 'commit 2' }}           |
      |         | git push --force-with-lease --force-if-includes |
    And the initial branches and lineage exist now
    And the initial commits exist now
