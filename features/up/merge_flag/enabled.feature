Feature: move up using the "merge" flag

  Scenario Outline: switching to child branch while merging open changes
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the current branch is "alpha"
    And an uncommitted file
    When I run "git-town up <FLAG>"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | alpha  | git checkout beta -m |
    And Git Town prints:
      """
        main
          alpha
      *     beta
      """

    Examples:
      | FLAG    |
      | --merge |
      | -m      |
