Feature: git town-hack: starting a new feature from the main branch (without remote repo)

  As a developer working on a new feature on the main branch and without a remote repository
  I want to be able to create a new up-to-date feature branch and continue my work there
  So that my work can exist on its own branch, code reviews remain effective, and my team productive.

  Background:
    Given my repo does not have a remote origin
    And the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE     | FILE NAME |
      | main   | local    | main_commit | main_file |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                     |
      | main        | git add -A                  |
      |             | git stash                   |
      |             | git branch new-feature main |
      |             | git checkout new-feature    |
      | new-feature | git stash pop               |
    And I am now on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH      | LOCATION | MESSAGE     | FILE NAME |
      | main        | local    | main_commit | main_file |
      | new-feature | local    | main_commit | main_file |
