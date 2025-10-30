Feature: handle conflicts between the current feature branch and the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git fetch --prune --tags                          |
      |         | git checkout main                                 |
      | main    | git -c rebase.updateRefs=false rebase origin/main |
      |         | git push                                          |
      |         | git checkout feature                              |
      | feature | git merge --no-edit --ff main                     |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And a merge is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | feature | git merge --abort |
    And no merge is now in progress
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local, origin | conflicting main commit    | conflicting_file | main content    |
      | feature | local         | conflicting feature commit | conflicting_file | feature content |

  @messyoutput
  Scenario: undo through another sync invocation
    When I run "git-town sync" and enter into the dialog:
      | DIALOG              | KEYS    |
      | unfinished runstate | 3 enter |
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | feature | git merge --abort |
    And Git Town prints:
      """
      Handle unfinished command: undo
      """
    And no merge is now in progress
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local, origin | conflicting main commit    | conflicting_file | main content    |
      | feature | local         | conflicting feature commit | conflicting_file | feature content |

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And a merge is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git commit --no-edit                    |
      |         | git merge --no-edit --ff origin/feature |
      |         | git push                                |
    And no merge is now in progress
    And all branches are now synchronized
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

  Scenario: resolve resulting in no changes and continue
    When I resolve the conflict in "conflicting_file" with "feature content"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git commit --no-edit                    |
      |         | git merge --no-edit --ff origin/feature |
      |         | git push                                |
    And no merge is now in progress
    And all branches are now synchronized
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git push                                |
    And no merge is now in progress
    And all branches are now synchronized
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

  @messyoutput
  Scenario: resolve and continue+run another program
    Given I resolve the conflict in "conflicting_file" with "feature content"
    When I run "git-town compress" and enter into the dialog:
      | DIALOG              | KEYS    |
      | unfinished runstate | 5 enter |
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git commit --no-edit                            |
      |         | git merge --no-edit --ff origin/feature         |
      |         | git push                                        |
      |         | git fetch --prune --tags                        |
      |         | git reset --soft main --                        |
      |         | git commit -m "conflicting feature commit"      |
      |         | git push --force-with-lease --force-if-includes |
    And Git Town prints:
      """
      Handle unfinished command: both
      """
    And no merge is now in progress
    And all branches are now synchronized
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |
