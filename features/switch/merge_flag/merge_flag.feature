@messyoutput
Feature: switch branches using the "merge" flag

  Scenario Outline: switching to another branch while merging open changes
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | current | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | other  | local    | other commit |
    And the current branch is "current"
    When I run "git-town switch <FLAG>" and enter into the dialogs:
      | DIALOG        | KEYS       |
      | switch-branch | down enter |
    Then Git Town runs the commands
      | BRANCH  | COMMAND               |
      | current | git checkout other -m |

    Examples:
      | FLAG    |
      | --merge |
      | -m      |
