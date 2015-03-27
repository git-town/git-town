Feature: git-sync-fork on a feature branch without open changes

  Background:
    Given I have a feature branch named "feature"
    And my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         | FILE NAME     |
      | main   | upstream | upstream commit | upstream_file |
    And I am on the "feature" branch
    When I run `git sync-fork`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                  |
      | feature | git checkout main        |
      | main    | git fetch upstream       |
      |         | git rebase upstream/main |
      |         | git push                 |
      |         | git checkout feature     |
    And I am still on the "feature" branch
    And I have the following commits
      | BRANCH | LOCATION                    | MESSAGE         | FILE NAME     |
      | main   | local, remote, and upstream | upstream commit | upstream_file |
