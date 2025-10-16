Feature: detaching a branch in Spanish

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
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-2 | local    | commit 2a |
      | branch-2 | local    | commit 2b |
    And the current branch is "branch-2"
    When I run "git-town detach" with these environment variables
      | LANG | es_ES.UTF-8 |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-2 | git fetch --prune --tags                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
    And Git Town prints:
      """
      Rebase aplicado satisfactoriamente y actualizado refs/heads/branch-2.
      """
    And this lineage exists now
      """
      main
        branch-1
        branch-2
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo" with these environment variables
      | LANG | es_ES.UTF-8 |
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                |
      | branch-2 | git reset --hard {{ sha 'commit 2b' }} |
    And Git Town prints:
      """
      HEAD está ahora
      """
    And the initial lineage exists now
    And the initial commits exist now
