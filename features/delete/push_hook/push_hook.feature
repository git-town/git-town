Feature: undo deleting the current feature branch with disabled push-hook

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | current | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | other   | local, origin | other commit   |
    And the current branch is "current" and the previous branch is "other"

  Scenario: set to "false"
    Given Git setting "git-town.push-hook" is "false"
    When I run "git-town delete"
    And I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git branch current {{ sha 'current commit' }} |
      |        | git push --no-verify -u origin current        |
      |        | git checkout current                          |
    And the initial branches and lineage exist now
    And the initial commits exist now

  Scenario: set to "true"
    Given Git setting "git-town.push-hook" is "true"
    When I run "git-town delete"
    And I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git branch current {{ sha 'current commit' }} |
      |        | git push -u origin current                    |
      |        | git checkout current                          |
    And the initial branches and lineage exist now
    And the initial commits exist now
