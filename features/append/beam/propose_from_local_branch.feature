@messyoutput
Feature: beam commits and uncommitted changes from a local branch onto a new child branch and propose

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE     |
      | main     | origin   | main commit |
      | existing | local    | commit 1    |
      | existing | local    | commit 2    |
      | existing | local    | commit 3    |
      | existing | local    | commit 4    |
      | existing | local    | commit 5    |
    And the current branch is "existing"
    And the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    And an uncommitted file
    And I ran "git add ."
    When I run "git-town append new --beam --commit --message uncommitted --propose" and enter into the dialog:
      | DIALOG                 | KEYS                             |
      | select commits 2 and 4 | down space down down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                       |
      | existing | git checkout -b new                                                                                           |
      | new      | git commit -m uncommitted                                                                                     |
      |          | git checkout existing                                                                                         |
      | existing | git -c rebase.updateRefs=false rebase --onto {{ sha-before-run 'commit 4' }}^ {{ sha-before-run 'commit 4' }} |
      |          | git -c rebase.updateRefs=false rebase --onto {{ sha-before-run 'commit 2' }}^ {{ sha-before-run 'commit 2' }} |
      |          | git checkout new                                                                                              |
      | new      | git -c rebase.updateRefs=false rebase existing                                                                |
      |          | git push -u origin new                                                                                        |
      |          | open https://github.com/git-town/git-town/compare/existing...new?expand=1&title=uncommitted                   |
      | new      | git checkout existing                                                                                         |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | origin        | main commit |
      | existing | local         | commit 1    |
      |          |               | commit 3    |
      |          |               | commit 5    |
      | new      | local, origin | commit 2    |
      |          |               | commit 4    |
      |          |               | uncommitted |
      |          | origin        | commit 1    |
      |          |               | commit 3    |
      |          |               | commit 5    |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                          |
      | existing | git reset --hard {{ sha-before-run 'commit 5' }} |
      |          | git branch -D new                                |
      |          | git push origin :new                             |
    And the initial commits exist now
    And the initial branches and lineage exist now
