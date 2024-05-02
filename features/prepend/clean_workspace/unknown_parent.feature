Feature: ask for missing parent information

  Scenario:
    Given the current branch is "old"
    When I run "git-town prepend new" and enter into the dialog:
      | DIALOG               | KEYS  |
      | parent branch of old | enter |
    Then it runs the commands
      | BRANCH | COMMAND                       |
      | old    | git fetch --prune --tags      |
      |        | git checkout main             |
      | main   | git rebase origin/main        |
      |        | git checkout old              |
      | old    | git merge --no-edit --ff main |
      |        | git push -u origin old        |
      |        | git checkout -b new main      |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |
