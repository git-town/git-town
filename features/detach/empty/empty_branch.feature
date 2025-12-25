Feature: detaching an empty branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-1 | local    | commit 1a |
      | branch-1 | local    | commit 1b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-2 | feature | branch-1 | local     |
      | branch-3 | feature | branch-2 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-3 | local    | commit 3a |
      | branch-3 | local    | commit 3b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-4 | feature | branch-3 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-4 | local    | commit 4a |
      | branch-4 | local    | commit 4b |
    And the current branch is "branch-2"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-2 | git fetch --prune --tags                                       |
      |          | git checkout branch-3                                          |
      | branch-3 | git -c rebase.updateRefs=false rebase --onto branch-1 branch-2 |
      |          | git checkout branch-4                                          |
      | branch-4 | git -c rebase.updateRefs=false rebase --onto branch-3 branch-2 |
      |          | git checkout branch-2                                          |
      | branch-2 | git -c rebase.updateRefs=false rebase --onto main branch-1     |
    And this lineage exists now
      """
      main
        branch-1
          branch-3
            branch-4
        branch-2
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                |
      | branch-2 | git reset --hard {{ sha 'commit 1b' }} |
    And the initial lineage exists now
    And the initial commits exist now
