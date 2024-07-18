Feature: syncs all feature branches (in a local repo)

  Background:
    Given a local Git repo clone
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | alpha | feature | main   | local     |
      | beta  | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | main   | local    | main commit  |
      | alpha  | local    | alpha commit |
      | beta   | local    | beta commit  |
    And the current branch is "alpha"
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                       |
      | alpha  | git merge --no-edit --ff main |
      |        | git checkout beta             |
      | beta   | git merge --no-edit --ff main |
      |        | git checkout alpha            |
    And the current branch is still "alpha"
    And all branches are now synchronized
