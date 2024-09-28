@messyoutput
Feature: switch to a new remote branch

  Scenario Outline:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | local-1  | feature | main   | local, origin |
      | local-2  | feature | main   | local, origin |
      | remote-1 | feature | main   | origin        |
    And the current branch is "local-2"
    And an uncommitted file
    And I ran "git fetch"
    When I run "git-town switch <FLAG>" and enter into the dialogs:
      | KEYS            |
      | down down enter |
    Then it runs the commands
      | BRANCH  | COMMAND               |
      | local-2 | git checkout remote-1 |
    And the current branch is now "remote-1"

    Examples:
      | FLAG  |
      | --all |
      | -a    |
