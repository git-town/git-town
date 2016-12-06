Feature: git town-hack: starting a new feature from a subfolder on the main branch (with remote repo)

  As a developer working in a subfolder on the main branch
  I want to be able to extract my open changes into a feature branch
  So that I can get them reviewed.


  Background:
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE       | FILE NAME        |
      | main   | local    | folder commit | new_folder/file1 |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git town-hack new-feature` in the "new_folder" folder


  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                           |
      | main        | git fetch --prune                 |
      | <none>      | cd <%= git_root_folder %>         |
      | main        | git add -A                        |
      |             | git stash                         |
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
      | main        | local and remote | folder commit |
      | new-feature | local and remote | folder commit |
