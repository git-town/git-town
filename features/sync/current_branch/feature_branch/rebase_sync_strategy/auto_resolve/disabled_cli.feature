@this
Feature: don't auto-resolve phantom merge conflicts

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                     | FILE NAME        | FILE CONTENT     |
      | main     | local    | conflicting main commit     | conflicting_file | main content     |
      | branch-1 | local    | commit 1                    | other_file       | content          |
      | branch-2 | local    | conflicting branch-2 commit | conflicting_file | branch-2 content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin deletes the "branch-1" branch
    And the current branch is "branch-2"
    When I run "git-town sync --auto-resolve=0"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-2 | git fetch --prune --tags                                   |
      |          | git checkout main                                          |
      | main     | git -c rebase.updateRefs=false rebase origin/main          |
      |          | git push                                                   |
      |          | git checkout branch-2                                      |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
      |          | GIT_EDITOR=true git rebase --continue                      |
      |          | git push --force-with-lease                                |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And Git Town prints something like:
      """
      could not apply .* conflicting branch-2 commit
      """
    And a rebase is now in progress
