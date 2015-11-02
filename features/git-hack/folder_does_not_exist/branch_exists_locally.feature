Feature: git hack: starting a new feature from the main branch (with remote repo)

  As a developer working on a new feature on the main branch
  I want to be able to create a new up-to-date feature branch and continue my work there
  So that my work can exist on its own branch, code reviews remain effective, and my team productive.


  Background:
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE       | FILE NAME        |
      | main   | remote   | main_commit   | main_file        |
      |        | local    | folder commit | new_folder/file1 |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git hack new-feature` in the "new_folder" folder


  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                           |
      | main        | git fetch --prune                 |
      | <none>      | cd <%= git_root_folder %>         |
      | main        | git stash -u                      |
      |             | git rebase origin/main            |
      |             | git push                          |
      |             | git checkout -b new-feature main  |
      | new-feature | git push -u origin new-feature    |
      |             | git stash pop                     |
      | <none>      | cd <%= git_folder "new_folder" %> |
    And I am in the project root folder
    And I end up on the "new-feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH      | LOCATION         | MESSAGE       |
      | main        | local and remote | main_commit   |
      |             |                  | folder commit |
      | new-feature | local and remote | main_commit   |
      |             |                  | folder commit |
