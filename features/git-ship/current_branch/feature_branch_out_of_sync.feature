Feature: Git Ship: feature branch out of sync with remote

  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE       |
      | feature | remote   | remote commit |
      | feature | local    | local commit  |
    And I run `git ship -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The 'feature' branch is out of sync. Run 'git sync' to resolve."
    And I am still on the "feature" branch
