Feature: sync the current feature branch (without a tracking branch or remote repo)

  Background:
    Given my repo does not have a remote origin
    And my repo has a local feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE              |
      | main    | local    | local main commit    |
      | feature | local    | local feature commit |
    And I am on the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git merge --no-edit main |
    And I am still on the "feature" branch
    And my repo now has the commits
      | BRANCH  | LOCATION | MESSAGE                          |
      | main    | local    | local main commit                |
      | feature | local    | local feature commit             |
      |         |          | local main commit                |
      |         |          | Merge branch 'main' into feature |
    And all branches are now synchronized
