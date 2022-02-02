Feature: abort the ship via empty commit message

  Background:
    Given my repo has the feature branches "feature" and "other-feature"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | main    | local, remote | main commit    | main_file    | main content    |
      | feature | local         | feature commit | feature_file | feature content |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file with name "feature_file" and content "conflicting content"
    When I run "git-town ship feature" and enter an empty commit message

  @skipWindows
  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                                     |
      | other-feature | git fetch --prune --tags                    |
      |               | git add -A                                  |
      |               | git stash                                   |
      |               | git checkout main                           |
      | main          | git rebase origin/main                      |
      |               | git checkout feature                        |
      | feature       | git merge --no-edit origin/feature          |
      |               | git merge --no-edit main                    |
      |               | git checkout main                           |
      | main          | git merge --squash feature                  |
      |               | git commit                                  |
      |               | git reset --hard                            |
      |               | git checkout feature                        |
      | feature       | git reset --hard {{ sha 'feature commit' }} |
      |               | git checkout main                           |
      | main          | git checkout other-feature                  |
      | other-feature | git stash pop                               |
    And it prints the error:
      """
      aborted because commit exited with error
      """
    And I am still on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
    And Git Town is still aware of this branch hierarchy
      | BRANCH        | PARENT |
      | feature       | main   |
      | other-feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "other-feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | main    | local, remote | main commit    |
      | feature | local         | feature commit |
    And Git Town is still aware of this branch hierarchy
      | BRANCH        | PARENT |
      | feature       | main   |
      | other-feature | main   |
