Feature: Git Ship: errors when the branch diff is empty

  Background:
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | branch  | location | file name   | file content   |
      | main    | remote   | common_file | common content |
      | feature | local    | common_file | common content |
    And I am on the "feature" branch
    When I run `git ship -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'feature' has no shippable changes"
    And I am still on the "feature" branch
