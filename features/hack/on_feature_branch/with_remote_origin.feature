Feature: git town-hack: starting a new feature from a feature branch (with remote repo)

  As a developer working on something unrelated to my current feature branch
  I want to be able to create a new up-to-date feature branch and continue my work there
  So that my work can exist on its own branch, code reviews remain effective, and my team productive.

  Background:
    Given my repo has a feature branch named "existing-feature"
    And the following commits exist in my repo
      | BRANCH           | LOCATION | MESSAGE                 |
      | main             | remote   | main commit             |
      | existing-feature | local    | existing feature commit |
    And I am on the "existing-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH           | COMMAND                     |
      | existing-feature | git fetch --prune --tags    |
      |                  | git add -A                  |
      |                  | git stash                   |
      |                  | git checkout main           |
      | main             | git rebase origin/main      |
      |                  | git branch new-feature main |
      |                  | git checkout new-feature    |
      | new-feature      | git stash pop               |
    And I am now on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH           | LOCATION      | MESSAGE                 |
      | main             | local, remote | main commit             |
      | existing-feature | local         | existing feature commit |
      | new-feature      | local         | main commit             |
