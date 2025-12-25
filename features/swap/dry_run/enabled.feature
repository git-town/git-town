Feature: swapping a feature branch in dry-run mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      | branch-2 | local, origin | commit 2 |
      | branch-3 | local, origin | commit 3 |
    And the current branch is "branch-2"
    When I run "git-town swap --dry-run"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                            |
      | branch-2 | git fetch --prune --tags                                                           |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1                         |
      |          | git checkout branch-1                                                              |
      | branch-1 | git -c rebase.updateRefs=false rebase --onto branch-2 main                         |
      |          | git checkout branch-3                                                              |
      | branch-3 | git -c rebase.updateRefs=false rebase --onto branch-1 {{ sha-initial 'commit 2' }} |
      |          | git checkout branch-2                                                              |
    And the initial branches and lineage exist now
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And the initial commits exist now
