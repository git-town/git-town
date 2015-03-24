Feature: git sync: syncing the current non-feature branch (with open changes and without remote repo)

  As a developer syncing a non-feature branch and without a remote repository
  I want to be able update my ongoing work to include the latest finished features from the rest of the team
  So that our collaboration remains effective.


  Background:
    Given I have branches named "qa" and "production"
    And my repo does not have a remote origin
    And my non-feature branches are configured as "qa" and "production"
    And I am on the "qa" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | qa     | local    | local commit | local_file |
      | main   | local    | main commit  | main_file  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`


  Scenario: no conflict
    Then it runs the Git commands
      | BRANCH | COMMAND       |
      | qa     | git stash -u  |
      | qa     | git stash pop |
    And I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | qa     | local    | local commit | local_file |
      | main   | local    | main commit  | main_file  |
