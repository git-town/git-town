Feature: git kill: errors if supplied branch does not exist

  As a developer mistyping the branch name to remove
  I should get an error that the given branch does not exist
  So that I can delete the correct branch


  Background:
    Given I am on the "main" branch


  Scenario: with open changes
    When I run `git kill non-existing-feature`, it errors
    Given I have an uncommitted file with name: "uncommitted" and content: "stuff"
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "There is no branch named 'non-existing-feature'"
    And I end up on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: without open changes
    When I run `git kill non-existing-feature`, it errors
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "There is no branch named 'non-existing-feature'"
    And I end up on the "main" branch
