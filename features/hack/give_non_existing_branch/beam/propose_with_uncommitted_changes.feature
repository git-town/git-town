@messyoutput
Feature: beam a commit and uncommitted changes onto a new feature branch and propose

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
      | existing | local, origin | commit 4    |
    And the current branch is "existing"
    And the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    And an uncommitted file
    And I ran "git add ."
    When I run "git-town hack new --beam --commit --message uncommitted --propose" and enter into the dialog:
      | DIALOG                 | KEYS                             |
      | select commits 2 and 4 | down space down down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                       |
      | existing | git checkout -b new main                                                                                      |
      | new      | git commit -m uncommitted                                                                                     |
      |          | git cherry-pick {{ sha-before-run 'commit 2' }}                                                               |
      |          | git cherry-pick {{ sha-before-run 'commit 4' }}                                                               |
      |          | git checkout existing                                                                                         |
      | existing | git -c rebase.updateRefs=false rebase --onto {{ sha-before-run 'commit 4' }}^ {{ sha-before-run 'commit 4' }} |
      |          | git -c rebase.updateRefs=false rebase --onto {{ sha-before-run 'commit 2' }}^ {{ sha-before-run 'commit 2' }} |
      |          | git push --force-with-lease --force-if-includes                                                               |
      |          | git checkout new                                                                                              |
      | new      | git push -u origin new                                                                                        |
      |          | open https://github.com/git-town/git-town/compare/new?expand=1&title=uncommitted                              |
      |          | git checkout existing                                                                                         |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | origin        | main commit |
      | existing | local, origin | commit 1    |
      |          |               | commit 3    |
      | new      | local, origin | uncommitted |
      |          |               | commit 2    |
      |          |               | commit 4    |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                          |
      | existing | git reset --hard {{ sha-before-run 'commit 4' }} |
      |          | git push --force-with-lease --force-if-includes  |
      |          | git branch -D new                                |
      |          | git push origin :new                             |
    And the initial commits exist now
    And the initial branches and lineage exist now
