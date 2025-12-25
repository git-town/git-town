Feature: move down using the "merge" flag

  Scenario Outline: switching to parent branch while merging open changes
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the current branch is "beta"
    And an uncommitted file
    When I run "git-town down <FLAG>"
    Then Git Town runs the commands
      | BRANCH | COMMAND               |
      | beta   | git checkout alpha -m |
    And Git Town prints:
      """
        main
      *   alpha
            beta
      """

    Examples:
      | FLAG    |
      | --merge |
      | -m      |
