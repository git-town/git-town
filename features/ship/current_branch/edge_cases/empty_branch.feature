Feature: cannot ship an empty branch

  Background:
    Given my repo has a feature branch "empty-feature"
    And my repo contains the commits
      | BRANCH        | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main          | remote   | main commit    | common_file | common content |
      | empty-feature | local    | feature commit | common_file | common content |
    And I am on the "empty-feature" branch
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                                     |
      | empty-feature | git fetch --prune --tags                    |
      |               | git checkout main                           |
      | main          | git rebase origin/main                      |
      |               | git checkout empty-feature                  |
      | empty-feature | git merge --no-edit origin/empty-feature    |
      |               | git merge --no-edit main                    |
      |               | git reset --hard {{ sha 'feature commit' }} |
      |               | git checkout main                           |
      | main          | git checkout empty-feature                  |
    And it prints the error:
      """
      the branch "empty-feature" has no shippable changes
      """
    And I am still on the "empty-feature" branch
    And Git Town is still aware of this branch hierarchy
      | BRANCH        | PARENT |
      | empty-feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "empty-feature" branch
    And my repo now has the following commits
      | BRANCH        | LOCATION      | MESSAGE        |
      | main          | local, remote | main commit    |
      | empty-feature | local         | feature commit |
    And Git Town is still aware of this branch hierarchy
      | BRANCH        | PARENT |
      | empty-feature | main   |
