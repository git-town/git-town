@skipWindows
Feature: no TTY

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the current branch is "existing"
    And tool "open" is installed
    When I run "git-town propose" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                             |
      | existing | git fetch --prune --tags                                            |
      |          | git push -u origin existing                                         |
      |          | Finding proposal from existing into main ... none                   |
      |          | open https://github.com/git-town/git-town/compare/existing?expand=1 |

  Scenario: undo
    When I run "git-town undo" in a non-TTY shell
    Then Git Town runs the commands
      | BRANCH   | COMMAND                   |
      | existing | git push origin :existing |
