Feature: commit three down

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
      | branch-4 | feature | branch-3 | local, origin |
    And the current branch is "branch-4"
    And an uncommitted file "changes" with content "my changes"
    And I ran "git add changes"
    When I run "git-town commit --down=3 -m down-commit"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                           |
      | branch-4 | git checkout branch-1             |
      | branch-1 | git commit -m down-commit         |
      |          | git checkout branch-2             |
      | branch-2 | git merge --no-edit --ff branch-1 |
      |          | git checkout branch-3             |
      | branch-3 | git merge --no-edit --ff branch-2 |
      |          | git checkout branch-4             |
      | branch-4 | git merge --no-edit --ff branch-3 |
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE     | FILE NAME | FILE CONTENT |
      | branch-1 | local    | down-commit | changes   | my changes   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                     |
      | branch-4 | git checkout branch-1                       |
      | branch-1 | git reset --hard {{ sha 'initial commit' }} |
      |          | git checkout branch-2                       |
      | branch-2 | git reset --hard {{ sha 'initial commit' }} |
      |          | git checkout branch-3                       |
      | branch-3 | git reset --hard {{ sha 'initial commit' }} |
      |          | git checkout branch-4                       |
      | branch-4 | git reset --hard {{ sha 'initial commit' }} |
    And the initial branches and lineage exist now
    And the initial commits exist now
