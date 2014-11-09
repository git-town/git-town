Feature: Git Kill

  Scenario: Killing a feature branch without open changes
    Given I have a feature branch named "good-feature"
    And I am on the "stupid-feature" branch
    When I run `git kill`
    Then I end up on the "main" branch
    And the branch "stupid-feature" is deleted everywhere
    And the branch "good-feature" still exists


  Scenario: Killing a feature branch with open changes

  Scenario: Cannot kill the main branch

  Scenario: Cannot kill a non-feature branch
