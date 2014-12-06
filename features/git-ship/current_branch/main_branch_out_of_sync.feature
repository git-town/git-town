Feature: Git Ship: errors if main branch is out of sync

  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE       | FILE NAME        | FILE CONTENT   |
      | main    | remote   | remote commit | conflicting_file | remote content |
      | main    | local    | local commit  | conflicting_file | local content  |
    And I run `git ship -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The 'main' branch is out of sync. Run 'git sync' to resolve."
    And I am still on the "feature" branch
