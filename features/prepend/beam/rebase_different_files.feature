@messyoutput
Feature: prepend a branch to a feature branch using the "rebase" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  |
      | old    | local, origin | commit 1 |
      | old    | local, origin | commit 2 |
      | old    | local, origin | commit 3 |
      | old    | local, origin | commit 4 |
    And the current branch is "old"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town prepend parent --beam" and enter into the dialog:
      | DIALOG          | KEYS                             |
      | commits to beam | space down down down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                 |
      | old    | git checkout -b parent main                                                                             |
      | parent | git cherry-pick {{ sha-initial 'commit 1' }}                                                            |
      |        | git cherry-pick {{ sha-initial 'commit 4' }}                                                            |
      |        | git checkout old                                                                                        |
      | old    | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 1' }}^ {{ sha-initial 'commit 1' }} |
      |        | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 4' }}^ {{ sha-initial 'commit 4' }} |
      |        | git -c rebase.updateRefs=false rebase parent                                                            |
      |        | git push --force-with-lease --force-if-includes                                                         |
      |        | git checkout parent                                                                                     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | old    | local, origin | commit 2 |
      |        |               | commit 3 |
      |        | origin        | commit 1 |
      |        |               | commit 4 |
      | parent | local         | commit 1 |
      |        |               | commit 4 |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |
    When I run "git town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                      |
      | parent | git fetch --prune --tags                                                     |
      |        | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |        | git push -u origin parent                                                    |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | old    | local, origin | commit 2 |
      |        |               | commit 3 |
      | parent | local, origin | commit 1 |
      |        |               | commit 4 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | parent | git checkout old                                |
      | old    | git reset --hard {{ sha 'commit 4' }}           |
      |        | git push --force-with-lease --force-if-includes |
      |        | git branch -D parent                            |
    And the initial commits exist now
    And the initial lineage exists now

  Scenario: amend the beamed commit
    And I amend this commit
      | BRANCH | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT    |
      | parent | local    | commit 4b | file_4    | amended content |
    And the current branch is "old"
    When I run "git town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                      |
      | old    | git fetch --prune --tags                                                     |
      |        | git checkout parent                                                          |
      | parent | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |        | git push -u origin parent                                                    |
      |        | git checkout old                                                             |
      | old    | git -c rebase.updateRefs=false rebase --onto parent {{ sha 'commit 4' }}     |
      |        | git push --force-with-lease --force-if-includes                              |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE   |
      | old    | local, origin | commit 2  |
      |        |               | commit 3  |
      | parent | local, origin | commit 1  |
      |        |               | commit 4b |
