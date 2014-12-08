Feature: Git Ship: feature branch out of sync with remote without open changes

  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE       |
      | feature | remote   | remote commit |
      | feature | local    | local commit  |
    And I am on the "other_feature" branch
    When I run `git ship feature -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The 'feature' branch is out of sync. Run 'git sync' to resolve."
    And I am still on the "other_feature" branch
