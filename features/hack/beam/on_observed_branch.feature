@messyoutput
Feature: beam from an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT | LOCATIONS     |
      | branch-1 | observed | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      |          | local         | commit 2 |
    And the current branch is "branch-1"
    When I run "git-town hack --beam branch-2" and enter into the dialog:
      | DIALOG          | KEYS        |
      | commits to beam | space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                 |
      | branch-1 | git checkout -b branch-2 main                                                                           |
      | branch-2 | git cherry-pick {{ sha-initial 'commit 1' }}                                                            |
      |          | git checkout branch-1                                                                                   |
      | branch-1 | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 1' }}^ {{ sha-initial 'commit 1' }} |
      |          | git push --force-with-lease --force-if-includes                                                         |
      |          | git checkout branch-2                                                                                   |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        branch-1
        branch-2
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 2 |
      | branch-2 | local         | commit 1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                  |
      | branch-2 | git checkout branch-1                                                    |
      | branch-1 | git reset --hard {{ sha 'commit 2' }}                                    |
      |          | git push --force-with-lease origin {{ sha-initial 'commit 1' }}:branch-1 |
      |          | git branch -D branch-2                                                   |
    And the initial branches and lineage exist now
    And the initial commits exist now
