Feature: Git Ship: errors when the branch diff is empty with open changes

  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION         | FILE NAME   | FILE CONTENT   |
      | main    | local and remote | common_file | common content |
      | feature | local            | common_file | common content |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship feature` while allowing errors


  Scenario: result
    Then I get the error "The branch 'feature' has no shippable changes"
    And I am still on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
