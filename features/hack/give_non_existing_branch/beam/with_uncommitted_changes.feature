@messyoutput
Feature: beam a commit and uncommitted changes onto a new feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | origin        | main commit |
      | existing | local, origin | commit 1    |
      | existing | local, origin | commit 2    |
      | existing | local, origin | commit 3    |
    And the current branch is "existing"
    And an uncommitted file
    And I ran "git add ."
    When I run "git-town hack new --beam --commit --message uncommitted" and enter into the dialog:
      | DIALOG          | KEYS             |
      | select commit 2 | down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                       |
      | existing | git checkout -b new main                                                                                      |
      | new      | git commit -m uncommitted                                                                                     |
      |          | git cherry-pick {{ sha-before-run 'commit 2' }}                                                               |
      |          | git checkout existing                                                                                         |
      | existing | git -c rebase.updateRefs=false rebase --onto {{ sha-before-run 'commit 2' }}^ {{ sha-before-run 'commit 2' }} |
      |          | git push --force-with-lease --force-if-includes                                                               |
      |          | git checkout existing                                                                                         |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | origin        | main commit |
      | existing | local, origin | commit 1    |
      |          |               | commit 3    |
      | new      | local         | uncommitted |
      |          |               | commit 2    |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                          |
      | existing | git reset --hard {{ sha-before-run 'commit 3' }} |
      |          | git push --force-with-lease --force-if-includes  |
      |          | git branch -D new                                |
    And the initial commits exist now
    And the initial branches and lineage exist now
