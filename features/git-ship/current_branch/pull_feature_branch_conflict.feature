Feature: Git Ship: handling conflicting remote feature branch updates when shipping the current feature branch


  Background:
    Given I am on the "feature-with-remote-conflicts" branch
    And the following commits exist in my repository
      | BRANCH                        | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature-with-remote-conflicts | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
      |                               | local    | local conflicting commit  | conflicting_file | local conflicting content  |
    And I run `git ship` while allowing errors


  Scenario: result
    Then I am still on the "feature-with-remote-conflicts" branch
    And my repo has a merge in progress
    And there is an abort script for "git ship"


  Scenario: aborting
    When I run `git ship --abort`
    Then I am still on the "feature-with-remote-conflicts" branch
    And there is no merge in progress
    And there is no abort script for "git ship" anymore
    And I still have the following commits
      | BRANCH                        | LOCATION | MESSAGE                   | FILES            |
      | feature-with-remote-conflicts | local    | local conflicting commit  | conflicting_file |
      |                               | remote   | remote conflicting commit | conflicting_file |
    And I still have the following committed files
      | BRANCH                        | FILES            | CONTENT                   |
      | feature-with-remote-conflicts | conflicting_file | local conflicting content |

