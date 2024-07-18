@messyoutput
Feature: ask for missing parent

  Scenario: on feature branch without parent
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town diff-parent" and enter into the dialog:
      | DIALOG                   | KEYS  |
      | parent branch of feature | enter |
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |
