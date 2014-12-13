Feature: git ship: don't ship empty feature branches

  As a developer trying to ship a feature branch that doesn't result in any changes on main
  I should be notified about this situation
  So that I can investigate and resolve this issue safely, and my users always see meaningful progress.


  Background:
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME   | FILE CONTENT   |
      | main    | remote   | common_file | common content |
      | feature | local    | common_file | common content |
    And I am on the "feature" branch
    When I run `git ship -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'feature' has no shippable changes"
    And I am still on the "feature" branch
