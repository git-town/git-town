Feature: git hack: starting a new feature from the main branch (with open changes and without remote repo)

  As a developer working on a new feature on the main branch and without a remote repository
  I want to be able to create a new up-to-date feature branch and continue my work there
  So that my work can exist on its own branch, code reviews remain effective, and my team productive.


  Background:
    Given my repo does not have a remote origin
    And the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE     | FILE NAME |
      | main   | local    | main_commit | main_file |
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack new_feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH      | COMMAND                          |
      | main        | git stash -u                     |
      | main        | git checkout -b new_feature main |
      | new_feature | git stash pop                    |
    And I end up on the "new_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH      | LOCATION | MESSAGE     | FILE NAME |
      | main        | local    | main_commit | main_file |
      | new_feature | local    | main_commit | main_file |
