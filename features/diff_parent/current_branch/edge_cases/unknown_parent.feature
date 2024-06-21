@messy
Feature: ask for missing parent

  Scenario: on feature branch without parent
    Given the current branch is "feature"
    When I run "git-town diff-parent" and enter into the dialog:
      | DIALOG                   | KEYS  |
      | parent branch of feature | enter |
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |
