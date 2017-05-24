Feature: git town-hack: starting a new feature from the main branch (with remote repo)

  As a developer working on a new feature on the main branch
  I want to be able to create a new up-to-date feature branch and continue my work there
  So that my work can exist on its own branch, code reviews remain effective, and my team productive.


  Background:
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE     |
      | main   | remote   | main_commit |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git-town hack new-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                          |
      | main        | git fetch --prune                |
      |             | git add -A                       |
      |             | git stash                        |
      |             | git rebase origin/main           |
      |             | git checkout -b new-feature main |
      | new-feature | git push -u origin new-feature   |
      |             | git stash pop                    |
    And I end up on the "new-feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH      | LOCATION         | MESSAGE     |
      | main        | local and remote | main_commit |
      | new-feature | local and remote | main_commit |
