Feature: git town-ship: errors when trying to ship the supplied feature branch that has no differences with the main branch


  Background:
    Given my repo has the feature branches "empty-feature" and "other-feature"
    And the following commits exist in my repo
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
