Feature: git town-hack: starting a new feature from the main branch (with remote repo)

  As a developer working on a new feature on the main branch
  I want to be able to create a new up-to-date feature branch and continue my work there
  So that my work can exist on its own branch, code reviews remain effective, and my team productive.

  Background:
    Given the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE     |
      | main   | remote   | main_commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                     |
      | main        | git fetch --prune --tags    |
      |             | git add -A                  |
      |             | git stash                   |
      |             | git rebase origin/main      |
      |             | git branch new-feature main |
      |             | git checkout new-feature    |
      | new-feature | git stash pop               |
    And I am now on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH      | LOCATION      | MESSAGE     |
      | main        | local, remote | main_commit |
      | new-feature | local         | main_commit |
