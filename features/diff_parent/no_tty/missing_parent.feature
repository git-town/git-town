Feature: no TTY, missing parent branch

  @this
  Scenario: feature branch
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE   | PARENT | LOCATIONS |
      | feature | (none) |        | local     |
    And the current branch is "feature"
    When I run "git-town diff-parent" in a non-TTY shell
    Then Git Town runs the commands
      | BRANCH | TYPE | COMMAND |
