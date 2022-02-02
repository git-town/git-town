Feature: syncing the current feature branch (without a tracking branch or remote repo)

  Background:
    Given my repo does not have a remote origin
    And my repo has a local feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE              |
      | main    | local    | local main commit    |
      | feature | local    | local feature commit |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git add -A               |
      |         | git stash                |
      |         | git merge --no-edit main |
      |         | git stash pop            |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
