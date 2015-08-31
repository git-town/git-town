Feature: git hack: starting a new feature from the main branch (with remote repo)

  As a developer working on a new feature on the main branch
  I want to be able to create a new up-to-date feature branch and continue my work there
  So that my work can exist on its own branch, code reviews remain effective, and my team productive.


  Background:
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE     | FILE NAME |
      | main   | remote   | main_commit | main_file |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git hack new-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                          |
      | main        | git fetch --prune                |
      |             | git stash -u                     |
      |             | git rebase origin/main           |
      |             | git checkout -b new-feature main |
      | new-feature | git stash pop                    |
    And I end up on the "new-feature" branch
    And I still have my uncommitted file
    And the branch "new-feature" has not been pushed to the repository
    And I have the following commits
      | BRANCH      | LOCATION         | MESSAGE     | FILE NAME |
      | main        | local and remote | main_commit | main_file |
      | new-feature | local            | main_commit | main_file |
