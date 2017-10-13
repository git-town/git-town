Feature: git town-rename-branch: errors when the destination branch exists locally

  As a developer trying to rename a branch to an already existing branch
  I want the command to abort with an error message
  So that I don't lose work by accidentally overwriting existing branches.


  Background:
    Given my repository has the feature branches "current-feature" and "existing-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION         | MESSAGE                 |
      | current-feature  | local and remote | current-feature commit  |
      | existing-feature | local and remote | existing-feature commit |
    And I am on the "current-feature" branch
    And my workspace has an uncommitted file
    When I run `git-town rename-branch current-feature existing-feature`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH          | COMMAND           |
      | current-feature | git fetch --prune |
    And it prints the error "A branch named 'existing-feature' already exists"
    And I am still on the "current-feature" branch
    And my workspace still contains my uncommitted file
    And my repository is left with my original commits
