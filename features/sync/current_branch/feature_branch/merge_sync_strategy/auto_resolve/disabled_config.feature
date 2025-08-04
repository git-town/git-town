Feature: disable auto-resolve phantom merge conflicts via configuration data

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE                     | FILE NAME        | FILE CONTENT |
      | branch-1 | local, origin | conflicting branch-1 commit | conflicting_file | content 1    |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                     | FILE NAME        | FILE CONTENT |
      | branch-2 | local    | conflicting branch-2 commit | conflicting_file | content 2    |
    And Git setting "git-town.sync-feature-strategy" is "merge"
    And Git setting "git-town.auto-resolve" is "no"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync"

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

  Scenario: undo
    When I run "git town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                             |
      | branch-2 | git merge --abort                                                   |
      |          | git checkout main                                                   |
      | main     | git reset --hard {{ sha 'initial commit' }}                         |
      |          | git branch branch-1 {{ sha-initial 'conflicting branch-1 commit' }} |
      |          | git checkout branch-2                                               |
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE                     | FILE NAME        | FILE CONTENT |
      | main     | origin   | conflicting branch-1 commit | conflicting_file | content 1    |
      | branch-1 | local    | conflicting branch-1 commit | conflicting_file | content 1    |
      | branch-2 | local    | conflicting branch-2 commit | conflicting_file | content 2    |
    And no merge is now in progress

  Scenario: run without resolving the conflicts
    When I run "git town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And a merge is still in progress

  Scenario: resolve the conflicts and continue
    When I resolve the conflict in "conflicting_file" with "content_2"
    And I run "git town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | branch-2 | git commit --no-edit                     |
      |          | git merge --no-edit --ff origin/branch-2 |
      |          | git push                                 |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                           | FILE NAME        | FILE CONTENT |
      | main     | local, origin | conflicting branch-1 commit       | conflicting_file | content 1    |
      | branch-2 | local, origin | conflicting branch-2 commit       | conflicting_file | content 2    |
      |          |               | Merge branch 'main' into branch-2 | conflicting_file | content_2    |
    And no merge is now in progress

  Scenario: resolve the conflicts, commit, and continue
    When I resolve the conflict in "conflicting_file" with "content_2"
    And I ran "git commit --no-edit"
    And I run "git town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | branch-2 | git merge --no-edit --ff origin/branch-2 |
      |          | git push                                 |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                           | FILE NAME        | FILE CONTENT |
      | main     | local, origin | conflicting branch-1 commit       | conflicting_file | content 1    |
      | branch-2 | local, origin | conflicting branch-2 commit       | conflicting_file | content 2    |
      |          |               | Merge branch 'main' into branch-2 | conflicting_file | content_2    |
    And no merge is now in progress
