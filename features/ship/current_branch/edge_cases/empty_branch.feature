Feature: does not ship an empty branch

  Background:
    Given the current branch is a feature branch "empty"
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME   | FILE CONTENT   |
      | main   | origin   | main commit  | common_file | common content |
      | empty  | local    | empty commit | common_file | common content |
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                   |
      | empty  | git fetch --prune --tags                  |
      |        | git checkout main                         |
      | main   | git rebase origin/main                    |
      |        | git checkout empty                        |
      | empty  | git merge --no-edit origin/empty          |
      |        | git merge --no-edit main                  |
      |        | git reset --hard {{ sha 'empty commit' }} |
      |        | git checkout main                         |
      | main   | git checkout empty                        |
    And it prints the error:
      """
      the branch "empty" has no shippable changes
      """
    And the current branch is still "empty"
    And the initial branch hierarchy exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And the current branch is still "empty"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | main commit  |
      | empty  | local         | empty commit |
    And the initial branch hierarchy exists
