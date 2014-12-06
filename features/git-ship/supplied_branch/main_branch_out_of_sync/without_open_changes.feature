Feature: Git Ship: errors if main branch is out of sync without open changes

  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE       | FILE NAME        | FILE CONTENT   |
      | main    | remote   | remote commit | conflicting_file | remote content |
      | main    | local    | local commit  | conflicting_file | local content  |
    And I am on the "other_feature" branch
    When I run `git ship feature -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The 'main' branch is out of sync. Run 'git sync' to resolve."
    And I am still on the "other_feature" branch
