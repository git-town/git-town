Feature: git sync: syncing the current non-feature branch (without open changes or remote repo)

  (see ./with_open_changes.feature)


  Background:
    Given I have branches named "qa" and "production"
    And my repo does not have a remote origin
    And my non-feature branches are configured as "qa" and "production"
    And I am on the "qa" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | qa     | local    | local commit | local_file |
      | main   | local    | main commit  | main_file  |
    When I run `git sync`


  Scenario: no conflict
    Then I am still on the "qa" branch
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | qa     | local    | local commit | local_file |
      | main   | local    | main commit  | main_file  |
