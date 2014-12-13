Feature: git ship: does not ship feature branches without changes

  As a developer shipping a feature branch that doesn't result in any changes on main
  I want to be notified about this situation
  So that I can investigate and resolve it safely, don't ship the wrong things, and my users see meaningful progress.


  Background:
    Given I have a feature branch named "no-real-changes"
    And the following commit exists in my repository
      | BRANCH          | LOCATION | FILE NAME   | FILE CONTENT   |
      | main            | remote   | common_file | common content |
      | no-real-changes | local    | common_file | common content |
    And I am on the "no-real-changes" branch
    When I run `git ship -m 'no real changes done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'no-real-changes' has no shippable changes"
    And I am still on the "no-real-changes" branch
