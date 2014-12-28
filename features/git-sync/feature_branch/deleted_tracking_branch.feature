Feature: git-sync restores deleted tracking branch

  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME    |
      | feature | local and remote | feature commit | feature_file |
    And the "feature" branch gets deleted on the remote
    And I am on the "feature" branch
    When I run `git sync`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                    |
      | feature | git checkout main          |
      | main    | git fetch --prune          |
      | main    | git rebase origin/main     |
      | main    | git checkout feature       |
      | feature | git merge --no-edit main   |
      | feature | git push -u origin feature |
    And I am still on the "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE        | FILES        |
      | feature | local and remote | feature commit | feature_file |
