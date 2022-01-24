Feature: git-town sync: restores deleted tracking branch

  As a developer syncing a feature branch whose tracking branch has been deleted
  I want a new tracking branch to be created
  So that my work is safe in case my local copy gets lost.

  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME    |
      | feature | local, remote | feature commit | feature_file |
    And the "feature" branch gets deleted on the remote
    And I am on the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune --tags   |
      |         | git checkout main          |
      | main    | git rebase origin/main     |
      |         | git checkout feature       |
      | feature | git merge --no-edit main   |
      |         | git push -u origin feature |
    And I am still on the "feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME    |
      | feature | local, remote | feature commit | feature_file |
