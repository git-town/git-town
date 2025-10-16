Feature: swapping empty branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the current branch is "branch-2"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                  |
      | branch-2 | git fetch --prune --tags                                                                 |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1                               |
      |          | git checkout branch-1                                                                    |
      | branch-1 | git -c rebase.updateRefs=false rebase --onto branch-2 main                               |
      |          | git checkout branch-3                                                                    |
      | branch-3 | git -c rebase.updateRefs=false rebase --onto branch-1 {{ sha-initial 'initial commit' }} |
      |          | git checkout branch-2                                                                    |
    And this lineage exists now
      """
      main
        branch-2
          branch-1
            branch-3
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial lineage exists now
    And the initial commits exist now
