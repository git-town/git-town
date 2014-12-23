Feature: git ship: don't ship empty feature branches (with open changes)

  (see ../current_branch/no_diff.feature)


  Background:
    Given I have feature branches named "empty-feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH        | LOCATION | FILE NAME   | FILE CONTENT   |
      | main          | remote   | common_file | common content |
      | empty-feature | local    | common_file | common content |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship empty-feature` while allowing errors


  Scenario: result
    Then I get the error "The branch 'empty-feature' has no shippable changes"
    And I am still on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
