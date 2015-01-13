Feature: git-sync-fork on a feature branch with open changes

  Background:
    Given I have a feature branch named "feature"
    And my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         | FILE NAME     |
      | main   | upstream | upstream commit | upstream_file |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-fork`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                  |
      | feature | git stash -u             |
      | feature | git checkout main        |
      | main    | git fetch upstream       |
      | main    | git rebase upstream/main |
      | main    | git push                 |
      | main    | git checkout feature     |
      | feature | git stash pop            |
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH | LOCATION                    | MESSAGE         | FILE NAME     |
      | main   | local, remote, and upstream | upstream commit | upstream_file |
