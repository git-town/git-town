Feature: compatibility between different sync-feature-strategy settings when editing dependent changes

  Scenario: I use rebase and my coworker uses merge
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "feature"
    And a coworker clones the repository
    And the coworker fetches updates
    And the coworker sets the "sync-feature-strategy" to "merge"
    And the coworker is on the "feature" branch
    And the coworker sets the parent branch of "feature" as "main"
    #
    # I make a commit and sync
    Given I add this commit to the current branch:
      | MESSAGE         | FILE NAME | FILE CONTENT                  |
      | my first commit | file.txt  | line 1: my content 1\nline 2: |
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                  |
      | feature | local, origin | my first commit | file.txt  | line 1: my content 1\nline 2: |
    And all branches are now synchronized
    And no rebase is now in progress
    #
    # coworker makes a conflicting local commit concurrently with me and then syncs
    Given the coworker adds this commit to their current branch:
      | MESSAGE               | FILE NAME | FILE CONTENT                        |
      | coworker first commit | file.txt  | line 1:\nline 2: coworker content 1 |
    When the coworker runs "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git merge --no-edit --ff origin/feature |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file.txt
      """
    And the coworkers workspace now contains file "file.txt" with content:
      """
      <<<<<<< HEAD
      line 1:
      line 2: coworker content 1
      =======
      line 1: my content 1
      line 2:
      >>>>>>> origin/feature

      """
    When the coworker resolves the conflict in "file.txt" with:
      """
      line 1: my content 1
      line 2: coworker content 1
      """
    And the coworker runs "git town continue" and closes the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND              |
      | feature | git commit --no-edit |
      |         | git push             |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE                                                    | FILE NAME | FILE CONTENT                                     |
      | feature | local, coworker, origin | my first commit                                            | file.txt  | line 1: my content 1\nline 2:                    |
      |         | coworker, origin        | coworker first commit                                      | file.txt  | line 1:\nline 2: coworker content 1              |
      |         |                         | Merge remote-tracking branch 'origin/feature' into feature | file.txt  | line 1: my content 1\nline 2: coworker content 1 |
    And the coworkers workspace now contains file "file.txt" with content:
      """
      line 1: my content 1
      line 2: coworker content 1
      """
    #
    # I add a conflicting commit locally and then sync
    Given I add this commit to the current branch:
      | MESSAGE          | FILE NAME | FILE CONTENT                 |
      | my second commit | file.txt  | line 1: my content 2\nline 2 |
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
    And Git Town prints something like:
      """
      could not apply \S+ my second commit
      """
    And file "file.txt" now has content:
      """
      <<<<<<< HEAD
      line 1: my content 1
      line 2: coworker content 1
      =======
      line 1: my content 2
      line 2
      >>>>>>> {{ sha-short 'my second commit' }} (my second commit)
      """
    When I resolve the conflict in "file.txt" with:
      """
      line 1: my content 2
      line 2: coworker content 1
      """
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                      |
      | feature | GIT_EDITOR=true git rebase --continue                                        |
      |         | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file.txt
      """
    # Note: strange that Git tries to apply "my second commit" first, and after that "my first commit".
    # In this case, the correct resolution is to use the version from "my second commit",
    # which the user has already resolved above.
    And Git Town prints something like:
      """
      could not apply \S+ my first commit
      """
    And file "file.txt" now has content:
      """
      <<<<<<< HEAD
      line 1:
      line 2: coworker content 1
      =======
      line 1: my content 1
      line 2:
      >>>>>>> {{ sha-short 'my first commit' }} (my first commit)
      """
    And a rebase is now in progress
    When I resolve the conflict in "file.txt" with:
      """
      line 1: my content 1
      line 2: coworker content 1
      """
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | GIT_EDITOR=true git rebase --continue           |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE                                                    | FILE NAME | FILE CONTENT                                     |
      | feature | local, coworker, origin | coworker first commit                                      | file.txt  | line 1:\nline 2: coworker content 1              |
      |         | local, origin           | my first commit                                            | file.txt  | line 1: my content 1\nline 2: coworker content 1 |
      |         |                         | my second commit                                           | file.txt  | line 1: my content 2\nline 2: coworker content 1 |
      |         | coworker                | my first commit                                            | file.txt  | line 1: my content 1\nline 2:                    |
      |         |                         | Merge remote-tracking branch 'origin/feature' into feature | file.txt  | line 1: my content 1\nline 2: coworker content 1 |
    And no rebase is now in progress
