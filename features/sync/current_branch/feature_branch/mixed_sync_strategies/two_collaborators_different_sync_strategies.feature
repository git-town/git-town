Feature: compatibility between different sync-feature-strategy settings

  Scenario: I use rebase and my coworker uses merge
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    Given Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "feature"
    And a coworker clones the repository
    And the coworker fetches updates
    And the coworker sets the "sync-feature-strategy" to "merge"
    And the coworker is on the "feature" branch
    And the coworker sets the parent branch of "feature" as "main"
    #
    # I make a commit and sync
    Given I add this commit to the current branch:
      | MESSAGE         | FILE NAME | FILE CONTENT     |
      | my first commit | file.txt  | my first content |
    When I run "git town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT     |
      | feature | local, origin | my first commit | file.txt  | my first content |
    And no rebase is now in progress
    And all branches are now synchronized
    #
    # coworker makes a conflicting local commit concurrently with me and then syncs
    Given the coworker adds this commit to their current branch:
      | MESSAGE               | FILE NAME | FILE CONTENT           |
      | coworker first commit | file.txt  | coworker first content |
    When the coworker runs "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git merge --no-edit --ff origin/feature |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file.txt
      """
    And file "file.txt" now has content:
      """
      my first content
      """
    When the coworker resolves the conflict in "file.txt" with "my and coworker first content"
    And the coworker runs "git town continue" and closes the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND              |
      | feature | git commit --no-edit |
      |         | git push             |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE                                                    | FILE NAME | FILE CONTENT                  |
      | feature | local, coworker, origin | my first commit                                            | file.txt  | my first content              |
      |         | coworker, origin        | coworker first commit                                      | file.txt  | coworker first content        |
      |         |                         | Merge remote-tracking branch 'origin/feature' into feature | file.txt  | my and coworker first content |
    And the coworkers workspace now contains file "file.txt" with content "my and coworker first content"
    #
    # I add a conflicting commit locally and then sync
    Given I add this commit to the current branch:
      | MESSAGE          | FILE NAME | FILE CONTENT      |
      | my second commit | file.txt  | my second content |
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                              |
      | feature | git fetch --prune --tags                             |
      |         | git push --force-with-lease --force-if-includes      |
      |         | git -c rebase.updateRefs=false rebase origin/feature |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file.txt
      """
    And file "file.txt" now has content:
      """
      <<<<<<< HEAD
      my and coworker first content
      =======
      my second content
      >>>>>>> {{ sha-short 'my second commit' }} (my second commit)
      """
    When I resolve the conflict in "file.txt" with "my second and coworker first content"
    And I run "git town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                      |
      | feature | GIT_EDITOR=true git rebase --continue                                        |
      |         | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |         | git checkout --theirs file.txt                                               |
      |         | git add file.txt                                                             |
      |         | GIT_EDITOR=true git rebase --continue                                        |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file.txt
      """
    And Git Town prints something like:
      """
      could not apply \S+ my first commit
      """
    And file "file.txt" now has content:
      """
      <<<<<<< HEAD
      my first content
      =======
      my second and coworker first content
      >>>>>>> {{ sha-short 'my second commit' }} (my second commit)
      """
    And a rebase is now in progress
    When I resolve the conflict in "file.txt" with "my second and coworker first content"
    And I run "git town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | GIT_EDITOR=true git rebase --continue           |
      |         | git push --force-with-lease --force-if-includes |
    And no rebase is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE                                                    | FILE NAME | FILE CONTENT                         |
      | feature | local, coworker, origin | coworker first commit                                      | file.txt  | coworker first content               |
      |         | local, origin           | my first commit                                            | file.txt  | my first content                     |
      |         |                         | my second commit                                           | file.txt  | my second and coworker first content |
      |         | coworker                | my first commit                                            | file.txt  | my first content                     |
      |         |                         | Merge remote-tracking branch 'origin/feature' into feature | file.txt  | my and coworker first content        |
