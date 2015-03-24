Feature: git sync: syncing the main branch (without open changes or remote repo)

  (see ./with_open_changes.feature)

  Background:
    Given I am on the "main" branch
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | LOCATION | MESSAGE      | FILE NAME  |
      | local    | local commit | local_file |
    When I run `git sync`


  Scenario: result
    Then I am still on the "main" branch
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | local commit | local_file |
