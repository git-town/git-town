Feature: display debug statistics

  Scenario: select another branch
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"
    When I run "git-town set-parent --debug" and answer the prompts:
      | PROMPT                                      | ANSWER        |
      | Please specify the parent branch of 'child' | [DOWN][ENTER] |
    Then it prints:
      """
      Ran 11 shell commands.
      """
    And this branch lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |
