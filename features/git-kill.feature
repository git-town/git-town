Feature: Git Kill

  Background:
    Given I have a feature branch named "good-feature"


  Scenario: Killing a feature branch without open changes
    Given I am on the "stupid-feature" branch
    When I run `git kill`
    Then I end up on the "main" branch
    And the branch "stupid-feature" is deleted everywhere
    And the branch "good-feature" still exists


  Scenario: Killing a feature branch with open changes
    Given I am on the "stupid-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill`
    Then I end up on the "main" branch
    And the branch "stupid-feature" is deleted everywhere
    And the branch "good-feature" still exists
    And I don't have an uncommitted file with name: "uncommitted"

  Scenario: Does not kill the main branch
    Given I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill` while allowing errors
    Then I get the error "You can only kill feature branches"
    And I am still on the "main" branch
    And the branch "good-feature" still exists
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: Cannot kill a non-feature branch
