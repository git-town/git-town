Feature: does not merge contribution branches

  Scenario Outline: disallowed branches
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | parent  | feature | main   | local     |
      | current | <TYPE>  | parent | local     |
    And the current branch is "current"
    When I run "git-town merge"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot merge branch "current" because it has no parent
      """

    Examples:
      | TYPE         |
      | contribution |
      | observed     |
      | perennial    |

  @this
  Scenario Outline: allowed branches
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | parent  | feature | main   | local     |
      | current | <TYPE>  | parent | local     |
    And the current branch is "current"
    When I run "git-town merge"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                         |
      | current | git fetch --prune --tags        |
      |         | git merge --no-edit --ff parent |
      |         | git branch -D parent            |

    Examples:
      | TYPE      |
      | feature   |
      | parked    |
      | prototype |
