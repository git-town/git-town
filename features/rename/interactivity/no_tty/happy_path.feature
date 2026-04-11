@skipWindows
Feature: no TTY

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the current branch is "existing"
    When I run "git-town rename new" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                        |
      | existing | git branch --move existing new |
      |          | git checkout new               |

  Scenario: undo
    When I run "git-town undo" in a non-TTY shell
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                        |
      | new      | git branch existing {{ sha 'initial commit' }} |
      |          | git checkout existing                          |
      | existing | git branch -D new                              |
    And the initial branches and lineage exist now
