Feature: syncs all feature branches (without remote repo)

  Background:
    Given my repo does not have an origin
    And my repo has the local feature branches "alpha" and "beta"
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | main   | local    | main commit  |
      | alpha  | local    | alpha commit |
      | beta   | local    | beta commit  |
    And I am on the "alpha" branch
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git merge --no-edit main |
      |        | git checkout beta        |
      | beta   | git merge --no-edit main |
      |        | git checkout alpha       |
    And I am still on the "alpha" branch
    And all branches are now synchronized
