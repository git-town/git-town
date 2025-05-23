@messyoutput
Feature: prepend a branch to a feature branch using the "rebase" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | old    | local, origin | commit 1 | file      | content 1    |
      | old    | local, origin | commit 2 | file      | content 2    |
      | old    | local, origin | commit 3 | file      | content 3    |
    And the current branch is "old"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town prepend parent --beam" and enter into the dialog:
      | DIALOG          | KEYS             |
      | select commit 2 | down space enter |
    Then Git Town runs the commands
      | BRANCH | COMMAND                                      |
      | old    | git checkout -b parent main                  |
      | parent | git cherry-pick {{ sha-initial 'commit 2' }} |
    And Git Town prints the error:
      """
      CONFLICT (modify/delete): file deleted in HEAD and modified in
      """
    And file "file" now has content:
      """
      content 2
      """
    And wait 1 second to ensure new Git timestamps
    When I run "git add file"
    And I run "git town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                      |
      | parent | git cherry-pick --continue                   |
      |        | git checkout old                             |
      | old    | git -c rebase.updateRefs=false rebase parent |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And Git Town prints something like:
      """
      error: could not apply .* commit 1
      """
    And file "file" now has content:
      """
      <<<<<<< HEAD
      content 2
      =======
      content 1
      >>>>>>> {{ sha-short "commit 1" }} (commit 1)
      """
    When I resolve the conflict in "file" with "content 1"
    And I run "git town continue"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | old    | GIT_EDITOR=true git rebase --continue           |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout parent                             |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | old    | local, origin | commit 1 | file      | content 1    |
      |        |               | commit 2 | file      | content 2    |
      |        |               | commit 3 | file      | content 3    |
      |        | origin        | commit 2 | file      | content 2    |
      | parent | local         | commit 2 | file      | content 2    |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: first sync after prepend
    When I run "git town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                      |
      | parent | git fetch --prune --tags                                                     |
      |        | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |        | git push -u origin parent                                                    |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | old    | local, origin | commit 1 |
      |        |               | commit 2 |
      |        |               | commit 3 |
      | parent | local, origin | commit 2 |
    And no uncommitted files exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | parent | git checkout old                                |
      | old    | git reset --hard {{ sha 'commit 3' }}           |
      |        | git push --force-with-lease --force-if-includes |
      |        | git branch -D parent                            |
    And the initial commits exist now
    And the initial lineage exists now

  Scenario: sync and amend the beamed commit
    And wait 1 second to ensure new Git timestamps
    When I run "git town sync"
    And wait 1 second to ensure new Git timestamps
    And I amend this commit
      | BRANCH | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT    |
      | parent | local    | commit 2b | file      | amended content |
    And the current branch is "old"
    When I run "git town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                       |
      | old    | git fetch --prune --tags                                                                      |
      |        | git checkout parent                                                                           |
      | parent | git push --force-with-lease --force-if-includes                                               |
      |        | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }}                  |
      |        | git checkout old                                                                              |
      | old    | git -c rebase.updateRefs=false rebase --onto parent {{ sha-in-origin-before-run 'commit 2' }} |
      |        | git checkout --theirs file                                                                    |
      |        | git add file                                                                                  |
      |        | GIT_EDITOR=true git rebase --continue                                                         |
      |        | git push --force-with-lease --force-if-includes                                               |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT    |
      | old    | local, origin | commit 1  | file      | content 1       |
      |        |               | commit 2  | file      | content 2       |
      |        |               | commit 3  | file      | content 3       |
      | parent | local, origin | commit 2b | file      | amended content |
