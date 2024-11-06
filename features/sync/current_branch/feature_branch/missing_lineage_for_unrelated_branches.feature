Feature: do not ask for lineage of branches that don't need to get synced

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
      | feature-2 | (none)  | main   | local, origin |
    And the current branch is "feature-1"
    And the commits
      | BRANCH    | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature-1 | local    | conflicting local commit  | conflicting_file | local content  |
      |           | origin   | conflicting origin commit | conflicting_file | origin content |
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                   |
      | feature-1 | git fetch --prune --tags                  |
      |           | git checkout main                         |
      | main      | git rebase origin/main --no-update-refs   |
      |           | git checkout feature-1                    |
      | feature-1 | git merge --no-edit --ff main             |
      |           | git merge --no-edit --ff origin/feature-1 |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And the current branch is still "feature-1"
    And a merge is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH    | COMMAND              |
      | feature-1 | git commit --no-edit |
      |           | git push             |
    And all branches are now synchronized
