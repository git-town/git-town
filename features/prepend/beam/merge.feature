@messyoutput
Feature: prepend a branch to a feature branch using the "merge" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | old    | local, origin | commit 1 | file 1    | content 1    |
      | old    | local, origin | commit 2 | file 2    | content 2    |
      | old    | local, origin | commit 3 | file 3    | content 3    |
      | old    | local, origin | commit 4 | file 4    | content 4    |
    And the current branch is "old"
    And Git setting "git-town.sync-feature-strategy" is "merge"
    When I run "git-town prepend parent --beam" and enter into the dialog:
      | DIALOG                 | KEYS                             |
      | select commits 2 and 4 | down space down down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                 |
      | old    | git checkout -b parent main                                                                             |
      | parent | git cherry-pick {{ sha-initial 'commit 2' }}                                                            |
      |        | git cherry-pick {{ sha-initial 'commit 4' }}                                                            |
      |        | git checkout old                                                                                        |
      | old    | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 2' }}^ {{ sha-initial 'commit 2' }} |
      |        | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 4' }}^ {{ sha-initial 'commit 4' }} |
      |        | git merge --no-edit --ff parent                                                                         |
      |        | git push --force-with-lease --force-if-includes                                                         |
      |        | git checkout parent                                                                                     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        |
      | old    | local, origin | commit 1                       |
      |        |               | commit 2                       |
      |        |               | commit 3                       |
      |        |               | Merge branch 'parent' into old |
      |        | origin        | commit 2                       |
      |        |               | commit 4                       |
      | parent | local         | commit 2                       |
      |        |               | commit 4                       |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |
    When I run "git town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                   |
      | parent | git fetch --prune --tags  |
      |        | git push -u origin parent |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        |
      | old    | local, origin | commit 1                       |
      |        |               | commit 2                       |
      |        |               | commit 3                       |
      |        |               | Merge branch 'parent' into old |
      | parent | local, origin | commit 2                       |
      |        |               | commit 4                       |

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
      | BRANCH | COMMAND                         |
      | old    | git fetch --prune --tags        |
      |        | git checkout parent             |
      | parent | git push -u origin parent       |
      |        | git checkout old                |
      | old    | git merge --no-edit --ff parent |
      |        | git push                        |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        |
      | old    | local, origin | commit 1                       |
      |        |               | commit 2                       |
      |        |               | commit 3                       |
      |        |               | commit 4                       |
      |        |               | Merge branch 'parent' into old |
      |        |               | Merge branch 'parent' into old |
      | parent | local, origin | commit 2                       |
      |        |               | commit 4b                      |
