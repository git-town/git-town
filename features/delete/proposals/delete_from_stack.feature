Feature: delete the current feature branch from a stack and update proposals

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | BODY       |
      |  1 | alpha         | main          | alpha body |
      |  2 | beta          | alpha         | beta body  |
    And Git setting "git-town.proposals-show-lineage" is "cli"
    And the current branch is "alpha"
    When I run "git-town delete"

  @debug @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
      |        | git push origin :alpha   |
      |        | git checkout beta        |
      | beta   | git branch -D alpha      |
    And the proposals are now
      | ID | SOURCE BRANCH | TARGET BRANCH | BODY       |
      |  1 | alpha         | main          | alpha body |
      |  2 | beta          | alpha         | beta body  |
    And this lineage exists now
      """
      main
        beta
      """
    And the branches are now
      | REPOSITORY    | BRANCHES   |
      | local, origin | main, beta |
    And no uncommitted files exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git branch current {{ sha 'current commit' }} |
      |        | git push -u origin current                    |
      |        | git checkout current                          |
    And the initial branches and lineage exist now
    And the initial commits exist now
