Feature: syncs all feature branches (without remote repo)

  Background:
    Given my repo does not have a remote origin
    And my repo has the local feature branches "feature-1" and "feature-2"
    And the following commits exist in my repo
      | BRANCH    | LOCATION | MESSAGE          |
      | main      | local    | main commit      |
      | feature-1 | local    | feature-1 commit |
      | feature-2 | local    | feature-2 commit |
    And I am on the "feature-1" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | feature-1 | git add -A               |
      |           | git stash                |
      |           | git merge --no-edit main |
      |           | git checkout feature-2   |
      | feature-2 | git merge --no-edit main |
      |           | git checkout feature-1   |
      | feature-1 | git stash pop            |
    And I am still on the "feature-1" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
