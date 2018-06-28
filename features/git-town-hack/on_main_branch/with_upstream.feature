Feature: git-hack: on the main branch with a upstream remote

  Background:
    Given my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run `git-town hack new-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                     |
      | main        | git fetch --prune --tags    |
      |             | git add -A                  |
      |             | git stash                   |
      |             | git rebase origin/main      |
      |             | git fetch upstream main     |
      |             | git rebase upstream/main    |
      |             | git push                    |
      |             | git branch new-feature main |
      |             | git checkout new-feature    |
      | new-feature | git stash pop               |
    And I am still on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH      | LOCATION                    | MESSAGE         |
      | main        | local, remote, and upstream | upstream commit |
      | new-feature | local                       | upstream commit |
