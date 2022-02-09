Feature: does not ship an empty branch

  Background:
    Given a feature branch "empty-feature"
    And the commits
      | BRANCH        | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main          | origin   | main commit    | common_file | common content |
      | empty-feature | local    | feature commit | common_file | common content |
    And the current branch is "empty-feature"
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
    And the current branch is still "empty-feature"
    And Git Town is still aware of the initial branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And the current branch is still "empty-feature"
    And now these commits exist
      | BRANCH        | LOCATION      | MESSAGE        |
      | main          | local, origin | main commit    |
      | empty-feature | local         | feature commit |
    And Git Town is still aware of the initial branch hierarchy
