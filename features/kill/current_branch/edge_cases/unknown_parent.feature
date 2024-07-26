Feature: ask for missing parent branch information

  Scenario:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE   | LOCATIONS |
      | feature | (none) | local     |
    And the current branch is "feature"
    When I run "git-town kill feature" and enter into the dialog:
      | DIALOG                   | KEYS  |
      | parent branch of feature | enter |
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git branch -D feature    |
    And no lineage exists now
