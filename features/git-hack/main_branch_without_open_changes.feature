Feature: git hack: starting a new feature from the main branch

  As a developer starting a new feature while being on the main branch
  I want to be able to create a new feature branch and start my work on it with one command
  So that it is easy to keep the code quality high by getting my code reviewed before merging into main.


  Background:
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE     | FILE NAME |
      | main   | remote   | main_commit | main_file |
    And I am on the "main" branch
    When I run `git hack new_feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                          |
      | main   | git fetch --prune                |
      | main   | git rebase origin/main           |
      | main   | git checkout -b new_feature main |
    And I end up on the "new_feature" branch
    And the branch "new_feature" has not been pushed to the repository
    And I have the following commits
      | BRANCH      | LOCATION         | MESSAGE     | FILE NAME |
      | main        | local and remote | main_commit | main_file |
      | new_feature | local            | main_commit | main_file |
