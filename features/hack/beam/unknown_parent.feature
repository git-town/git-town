@messyoutput
Feature: beam commits onto a new feature branch when the parent branch is unknown

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE   | PARENT | LOCATIONS     |
      | existing | (none) |        | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE     |
      | main     | origin   | main commit |
      | existing | local    | commit 1    |
      | existing | local    | commit 2    |
    And the current branch is "existing"
    When I run "git-town hack new --beam" and enter into the dialog:
      | DIALOG                       | KEYS             |
      | parent branch for "existing" | enter            |
      | commits to beam              | down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                 |
      | existing | git checkout -b new main                                                                                |
      | new      | git cherry-pick {{ sha-initial 'commit 2' }}                                                            |
      |          | git checkout existing                                                                                   |
      | existing | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 2' }}^ {{ sha-initial 'commit 2' }} |
      |          | git push --force-with-lease --force-if-includes                                                         |
      |          | git checkout new                                                                                        |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        existing
        new
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | origin        | main commit |
      | existing | local, origin | commit 1    |
      | new      | local         | commit 2    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                |
      | new      | git checkout existing                                                  |
      | existing | git reset --hard {{ sha-initial 'commit 2' }}                          |
      |          | git push --force-with-lease origin {{ sha 'initial commit' }}:existing |
      |          | git branch -D new                                                      |
    And the initial branches and lineage exist now
    And the initial commits exist now
