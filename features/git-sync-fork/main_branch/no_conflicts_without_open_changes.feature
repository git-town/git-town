Feature: git-sync-fork on the main branch without open changes

  Background:
    Given my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         | FILE NAME     |
      | main   | upstream | upstream commit | upstream_file |
    And I am on the "main" branch
    When I run `git sync-fork`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                  |
      | main   | git fetch upstream       |
      |        | git rebase upstream/main |
      |        | git push                 |
    And I am still on the "main" branch
    And I have the following commits
      | BRANCH | LOCATION                    | MESSAGE         | FILE NAME     |
      | main   | local, remote, and upstream | upstream commit | upstream_file |
