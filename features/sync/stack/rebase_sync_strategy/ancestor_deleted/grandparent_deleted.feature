Feature: a grandparent was deleted remotely while the parent and the child remain

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | branch-1 | local, origin | branch-1 commit |
      | branch-2 | local, origin | branch-2 commit |
      | branch-3 | local, origin | branch-3 commit |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "branch-3"
    And I ran "git-town sync --stack"
    And origin deletes the "branch-1" branch
    When I run "git-town sync --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                      |
      | branch-3 | git fetch --prune --tags                                                                     |
      |          | git checkout branch-2                                                                        |
      | branch-2 | git pull                                                                                     |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1                                   |
      |          | git push --force-with-lease                                                                  |
      |          | git checkout branch-3                                                                        |
      | branch-3 | git -c rebase.updateRefs=false rebase --onto branch-2 {{ sha-before-run 'branch-2 commit' }} |
      |          | git push --force-with-lease --force-if-includes                                              |
      |          | git branch -D branch-1                                                                       |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        branch-2
          branch-3
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | branch-2 | local, origin | branch-2 commit |
      | branch-3 | local, origin | branch-3 commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-3 | git checkout branch-2                           |
      | branch-2 | git reset --hard {{ sha 'branch-2 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'branch-3 commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git branch branch-1 {{ sha 'branch-1 commit' }} |
    And the initial branches and lineage exist now
    And the initial commits exist now
