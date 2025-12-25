Feature: adding additional stash entries while resolving conflicts

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature | local    | conflicting local commit  | conflicting_file | local content  |
      |         | origin   | conflicting origin commit | conflicting_file | origin content |
    And the current branch is "feature"
    And an uncommitted file
    And I run "git-town sync"
    And Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git add -A                              |
      |         | git stash -m "Git Town WIP"             |
      |         | git merge --no-edit --ff origin/feature |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And a merge is now in progress
    And I resolve the conflict in "conflicting_file"
    And I run "git add ."
    And I run "git commit --no-edit"
    And I add an unrelated stash entry with file "stashed_file"

  Scenario: continue
    When I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                |
      | feature | git push               |
      |         | git stash pop          |
      |         | git restore --staged . |
    And no merge is now in progress
    And an uncommitted file "stashed_file" exists now
    And all branches are now synchronized

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                |
      | feature | git merge --abort      |
      |         | git stash pop          |
      |         | git restore --staged . |
    And no merge is now in progress
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                                                    | FILE NAME        | FILE CONTENT     |
      | feature | local         | conflicting local commit                                   | conflicting_file | local content    |
      |         | local, origin | conflicting origin commit                                  | conflicting_file | origin content   |
      |         | local         | Merge remote-tracking branch 'origin/feature' into feature | conflicting_file | resolved content |
    And an uncommitted file "stashed_file" exists now
