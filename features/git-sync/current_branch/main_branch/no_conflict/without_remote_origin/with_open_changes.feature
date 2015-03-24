Feature: git sync: syncing the main branch (with open changes and without remote repo)

  As a developer syncing the main branch and without a remote repository
  I want to be able update my ongoing work to include the latest finished features from the rest of the team
  So that our collaboration remains effective.


  Background:
    Given I am on the "main" branch
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | LOCATION | MESSAGE      | FILE NAME  |
      | local    | local commit | local_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND       |
      | main   | git stash -u  |
      | main   | git stash pop |
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | local commit | local_file |
