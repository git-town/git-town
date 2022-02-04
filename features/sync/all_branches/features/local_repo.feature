Feature: syncs all feature branches (without remote repo)

  Background:
    Given my repo does not have a remote origin
    And my repo has the local feature branches "one" and "two"
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main commit |
      | one    | local    | one commit  |
      | two    | local    | two commit  |
    And I am on the "one" branch
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | one    | git merge --no-edit main |
      |        | git checkout two         |
      | two    | git merge --no-edit main |
      |        | git checkout one         |
    And I am still on the "one" branch
    And all branches are now synchronized
