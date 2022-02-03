Feature: does not ship empty feature branches

  Background:
    Given my repo has the feature branches "empty-feature" and "other-feature"
    And my repo contains the commits
      | BRANCH        | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main          | remote   | main commit    | common_file | common content |
      | empty-feature | local    | feature commit | common_file | common content |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town ship empty-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                                     |
      | other-feature | git fetch --prune --tags                    |
      |               | git add -A                                  |
      |               | git stash                                   |
      |               | git checkout main                           |
      | main          | git rebase origin/main                      |
      |               | git checkout empty-feature                  |
      | empty-feature | git merge --no-edit origin/empty-feature    |
      |               | git merge --no-edit main                    |
      |               | git reset --hard {{ sha 'feature commit' }} |
      |               | git checkout main                           |
      | main          | git checkout other-feature                  |
      | other-feature | git stash pop                               |
    And it prints the error:
      """
      the branch "empty-feature" has no shippable changes
      """
    And I am still on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And Git Town still has the original branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "other-feature" branch
    And my repo now has the following commits
      | BRANCH        | LOCATION      | MESSAGE        |
      | main          | local, remote | main commit    |
      | empty-feature | local         | feature commit |
    And Git Town still has the original branch hierarchy
