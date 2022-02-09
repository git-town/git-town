Feature: sync the current feature branch (in a local repo)

  Background:
    Given my repo does not have an origin
    And my repo has a local feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | main    | local    | main commit    |
      | feature | local    | feature commit |
    And I am on the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git merge --no-edit main |
    And all branches are now synchronized
    And I am still on the "feature" branch
    And now these commits exist
      | BRANCH  | LOCATION | MESSAGE                          |
      | main    | local    | main commit                      |
      | feature | local    | feature commit                   |
      |         |          | main commit                      |
      |         |          | Merge branch 'main' into feature |
