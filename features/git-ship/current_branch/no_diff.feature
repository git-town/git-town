Feature: Git Ship: errors when the branch diff is empty

  Background:
    Given I have a feature branch named "empty-feature"
    And the following commit exists in my repository
      | BRANCH        | LOCATION | FILE NAME   | FILE CONTENT   |
      | main          | remote   | common_file | common content |
      | empty-feature | local    | common_file | common content |
    And I am on the "empty-feature" branch
    When I run `git ship -m 'no real changes done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'empty-feature' has no shippable changes"
    And I am still on the "empty-feature" branch
