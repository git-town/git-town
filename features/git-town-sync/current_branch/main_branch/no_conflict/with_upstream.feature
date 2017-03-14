Feature: git-sync: on the main branch with a upstream remote

  Background:
    Given my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `gt sync`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune        |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/main   |
      |        | git fetch upstream       |
      |        | git rebase upstream/main |
      |        | git push                 |
      |        | git push --tags          |
      |        | git stash pop            |
    And I am still on the "main" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH | LOCATION                    | MESSAGE         |
      | main   | local, remote, and upstream | upstream commit |
