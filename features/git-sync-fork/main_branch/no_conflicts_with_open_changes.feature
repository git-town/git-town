Feature: git-sync-fork on the main branch with open changes

  Background:
    Given my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         | FILE NAME     |
      | main   | upstream | upstream commit | upstream_file |
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-fork`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                  |
      | main   | git stash -u             |
      | main   | git fetch upstream       |
      | main   | git rebase upstream/main |
      | main   | git push                 |
      | main   | git stash pop            |
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH | LOCATION                    | MESSAGE         | FILE NAME     |
      | main   | local, remote, and upstream | upstream commit | upstream_file |
    And now I have the following committed files
      | BRANCH | FILES         |
      | main   | upstream_file |
