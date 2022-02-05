Feature: does not ship empty feature branches

  Background:
    Given my repo has the feature branches "empty" and "other"
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main   | remote   | main commit    | common_file | common content |
      | empty  | local    | feature commit | common_file | common content |
    And I am on the "other" branch
    And my workspace has an uncommitted file
    When I run "git-town ship empty"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | other  | git fetch --prune --tags                    |
      |        | git add -A                                  |
      |        | git stash                                   |
      |        | git checkout main                           |
      | main   | git rebase origin/main                      |
      |        | git checkout empty                          |
      | empty  | git merge --no-edit origin/empty            |
      |        | git merge --no-edit main                    |
      |        | git reset --hard {{ sha 'feature commit' }} |
      |        | git checkout main                           |
      | main   | git checkout other                          |
      | other  | git stash pop                               |
    And it prints the error:
      """
      the branch "empty" has no shippable changes
      """
    And I am still on the "other" branch
    And my workspace still contains my uncommitted file
    And Git Town still knows the initial branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "other" branch
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, remote | main commit    |
      | empty  | local         | feature commit |
    And Git Town still knows the initial branch hierarchy
