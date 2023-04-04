Feature: display debug statistics

  Scenario: feature branch
    And the current branch is a feature branch "feature"
    When I run "git-town diff-parent --debug"
    Then it prints:
      """
      Ran 9 shell commands.
      """
