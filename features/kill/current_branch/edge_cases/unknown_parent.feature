Feature: ask for missing parent branch information

  Scenario:
    Given a Git repo with origin
    And the branch
      | NAME   | TYPE   | LOCATIONS |
      | branch | (none) | local     |
    And the current branch is "branch"
    When I run "git-town kill"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | branch | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git branch -D branch     |
    And no lineage exists now
