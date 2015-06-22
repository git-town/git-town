Feature: git ship: errors when trying to ship the current feature branch that has no differences with the main branch

  As a developer shipping a branch that has no differences with the main branch
  I should see an error telling me about this
  So that I can investigate this issue, and my users always see meaningful progress.


  Background:
    Given I have a feature branch named "empty-feature"
    And the following commit exists in my repository
      | BRANCH        | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main          | remote   | main commit    | common_file | common content |
      | empty-feature | local    | feature commit | common_file | common content |
    And I am on the "empty-feature" branch
    When I run `git ship`


  Scenario: result
    Then it runs the Git commands
      | BRANCH        | COMMAND                                      |
      | empty-feature | git fetch --prune                            |
      |               | git checkout main                            |
      | main          | git rebase origin/main                       |
      |               | git checkout empty-feature                   |
      | empty-feature | git merge --no-edit origin/empty-feature     |
      |               | git merge --no-edit main                     |
      |               | git reset --hard <%= sha 'feature commit' %> |
      |               | git checkout main                            |
      | main          | git checkout empty-feature                   |
    And I get the error "The branch 'empty-feature' has no shippable changes"
    And I am still on the "empty-feature" branch
