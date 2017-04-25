Feature: git town-kill: errors if supplied branch does not exist

  As a developer mistyping the branch name to remove
  I should get an error that the given branch does not exist
  So that I can delete the correct branch


  Background:
    Given I am on the "main" branch


  Scenario: result
    When I run `git-town kill non-existing-feature`
    Given I have an uncommitted file
    Then it runs the commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "There is no branch named 'non-existing-feature'"
    And I end up on the "main" branch
    And I still have my uncommitted file
