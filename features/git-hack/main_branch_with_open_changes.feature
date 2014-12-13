Feature: git hack: moving existing open changes from the main branch into a new feature branch

  As a developer working on a new feature while still being on the main branch
  I want to be able to create a new feature branch and continue my work on it with one command
  So that it is easy to keep the code quality high by getting my code reviewed before merging into main.


  Background:
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE     | FILE NAME |
      | main   | remote   | main_commit | main_file |
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                      |
      | main    | git stash -u                 |
      | main    | git fetch --prune            |
      | main    | git rebase origin/main       |
      | main    | git checkout -b feature main |
      | feature | git stash pop                |
    And I end up on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the branch "feature" has not been pushed to the repository
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILES     |
      | main    | local and remote | main_commit | main_file |
      | feature | local            | main_commit | main_file |
    And now I have the following committed files
      | BRANCH  | FILES     |
      | main    | main_file |
      | feature | main_file |
