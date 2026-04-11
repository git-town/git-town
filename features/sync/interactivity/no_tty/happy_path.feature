@skipWindows
Feature: no TTY

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | local    | main commit     |
      | existing | local    | existing commit |
    And the current branch is "existing"
    When I run "git-town sync" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                       |
      | existing | git merge --no-edit --ff main |

  Scenario: undo
    When I run "git-town undo" in a non-TTY shell
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                      |
      | existing | git reset --hard {{ sha 'existing commit' }} |
