Feature: git sync: syncing the main branch (without remote repo)

  As a developer syncing the main branch and without a remote repository
  I want to be able update my ongoing work to include the latest finished features from the rest of the team
  So that our collaboration remains effective.


  Background:
    Given my repo does not have a remote origin
    And I am on the "main" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | local commit | local_file |
    And I have an uncommitted file
    When I run `git sync`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND       |
      | main   | git stash -a  |
      |        | git stash pop |
    And I am still on the "main" branch
    And I still have my uncommitted file
    And I am left with my original commits
