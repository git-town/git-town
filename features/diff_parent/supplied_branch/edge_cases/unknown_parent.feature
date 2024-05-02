Feature: ask for missing parent

  Scenario: branch without parent
    Given a branch "feature"
    And the current branch is "main"
    When I run "git-town diff-parent feature" and enter into the dialog:
      | DIALOG                   | KEYS  |
      | parent branch of feature | enter |
    Then it runs the commands
      | BRANCH | COMMAND                |
      | main   | git diff main..feature |
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |
