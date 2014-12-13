Feature: Git Ship: errors when the branch diff is empty without open changes

  Background:
    Given I have feature branches named "no-real-changes" and "other_feature"
    And the following commit exists in my repository
      | BRANCH          | LOCATION | FILE NAME   | FILE CONTENT   |
      | main            | remote   | common_file | common content |
      | no-real-changes | local    | common_file | common content |
    And I am on the "other_feature" branch
    When I run `git ship no-real-changes` while allowing errors


  Scenario: result
    Then I get the error "The branch 'no-real-changes' has no shippable changes"
    And I am still on the "other_feature" branch
