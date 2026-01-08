Feature: commit down

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
      | branch-1 | local, origin | commit 1a | file_1    | content 1    |
      | branch-2 | local, origin | commit 2a | file_2    | content 2    |
    And the current branch is "branch-2"
    And an uncommitted file "changes" with content "my changes"
    And I ran "git add changes"
    When I run "git-town commit --down -m commit-2b"

  @debug @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-2 | git fetch --prune --tags |
      |          | git checkout branch-1    |
      | branch-1 | git commit -m commit-2b  |
      |          | git checkout branch-2    |
    # And these commits exist now
    #   | BRANCH   | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
    #   | branch-1 | local, origin | commit 1a | file_1    | content 1    |
    #   | branch-2 | local, origin | commit 2a | file_2    | content 2    |
    #   | branch-2 | local, origin | commit 2b | file_2    | content 2    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git reset --hard {{ sha 'commit 1' }}           |
      |         | git push --force-with-lease --force-if-includes |
    And the initial branches and lineage exist now
    And the initial commits exist now
