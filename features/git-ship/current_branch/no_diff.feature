Feature: git ship: don't ship empty features

  As a developer trying to ship a feature branch that doesn't result in any changes on main
  I should see an error telling me about this
  So that I can investigate this issue, and my users always see meaningful progress.


  Background:
    Given I have a feature branch named "empty-feature"
    And the following commit exists in my repository
      | BRANCH        | LOCATION | FILE NAME   | FILE CONTENT   |
      | main          | remote   | common_file | common content |
      | empty-feature | local    | common_file | common content |
    And I am on the "empty-feature" branch
    When I run `git ship` while allowing errors


  Scenario: result
    Then I get the error "The branch 'empty-feature' has no shippable changes"
    And I am still on the "empty-feature" branch
