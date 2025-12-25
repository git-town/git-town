Feature: sync a stack that makes conflicting changes

  Scenario: all branches in the stack change the same file to different values and this hasn't been synced yet
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT  |
      | main   | origin        | main commit  | file      | main content  |
      | alpha  | local, origin | alpha commit | file      | alpha content |
      | beta   | local, origin | beta commit  | file      | beta content  |
    And the current branch is "alpha"
    When I run "git-town sync --all"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | alpha  | git fetch --prune --tags                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git checkout alpha                                |
      | alpha  | git merge --no-edit --ff main                     |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    When I resolve the conflict in "file" with "resolved alpha content"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                        |
      | alpha  | git commit --no-edit           |
      |        | git push                       |
      |        | git checkout beta              |
      | beta   | git merge --no-edit --ff alpha |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And a merge is now in progress
    When I resolve the conflict in "file" with "resolved beta content"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | beta   | git commit --no-edit |
      |        | git push             |
      |        | git checkout alpha   |
      | alpha  | git push --tags      |
    And no merge is now in progress
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        | FILE NAME | FILE CONTENT           |
      | main   | local, origin | main commit                    | file      | main content           |
      | alpha  | local, origin | alpha commit                   | file      | alpha content          |
      |        |               | Merge branch 'main' into alpha | file      | resolved alpha content |
      | beta   | local, origin | beta commit                    | file      | beta content           |
      |        |               | Merge branch 'alpha' into beta | file      | resolved beta content  |
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                               |
      | alpha  | git reset --hard {{ sha-initial 'alpha commit' }}     |
      |        | git push --force-with-lease --force-if-includes       |
      |        | git checkout beta                                     |
      | beta   | git reset --hard {{ sha-initial 'beta commit' }}      |
      |        | git push --force-with-lease --force-if-includes       |
      |        | git checkout main                                     |
      | main   | git reset --hard {{ sha-in-origin 'initial commit' }} |
      |        | git checkout alpha                                    |
    And the initial commits exist now
    And the initial branches and lineage exist now
