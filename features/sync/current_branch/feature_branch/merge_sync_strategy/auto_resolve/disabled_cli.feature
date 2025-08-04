Feature: disable auto-resolve phantom merge conflicts via CLI

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE                     | FILE NAME        | FILE CONTENT |
      | branch-1 | local, origin | conflicting branch-1 commit | conflicting_file | content 1    |
      | branch-2 | local         | conflicting branch-2 commit | conflicting_file | content 2    |
    And Git setting "git-town.sync-feature-strategy" is "merge"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync --auto-resolve=0"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                           |
      | branch-2 | git fetch --prune --tags                          |
      |          | git checkout main                                 |
      | main     | git -c rebase.updateRefs=false rebase origin/main |
      |          | git branch -D branch-1                            |
      |          | git checkout branch-2                             |
      | branch-2 | git merge --no-edit --ff main                     |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And a merge is now in progress

  @this
  Scenario: undo
    When I run "git town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                             |
      | branch-2 | git checkout main                                                   |
      |          | git branch branch-1 {{ sha-initial 'conflicting branch-1 commit' }} |
