Feature: handling conflicting remote main branch updates when hacking with open changes


  Background:
    Given I have a feature branch named "feature"
    Given the following commit exists in my repository
      | branch | location | message                   | file name        | file content   |
      | main   | remote   | remote_conflicting_commit | conflicting_file | remote content |
      | main   | local    | local_conflicting_commit  | conflicting_file | local content  |
    And I am on the "feature" branch
    When I run `git hack other_feature` while allowing errors


  Scenario: result
    Then my repo has a rebase in progress
    And there is an abort script for "git hack"


  Scenario: aborting
    When I run `git hack --abort`
    Then I end up on the "feature" branch
    And there is no rebase in progress
    And there is no abort script for "git hack" anymore
    And I have the following commits
      | branch | location | message                   | files            |
      | main   | remote   | remote_conflicting_commit | conflicting_file |
      | main   | local    | local_conflicting_commit  | conflicting_file |
