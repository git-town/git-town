Feature: git town-rename-branch: errors if the feature branch does not exist

  As a developer mistyping the feature branch name to rename
  I should get an error that the given branch does not exist
  So that I can rename the correct branch.


  Background:
    Given the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE     |
      | main   | local and remote | main commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run `git-town rename-branch non-existing-feature renamed-feature`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And it prints the error "There is no branch named 'non-existing-feature'"
    And I end up on the "main" branch
    And my workspace still contains my uncommitted file
    And my repository is left with my original commits
