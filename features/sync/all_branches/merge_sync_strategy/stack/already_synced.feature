Feature: sync a stack making independent changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | alpha  | local, origin | alpha 1 | alpha_1   | alpha 1      |
      |        |               | alpha 2 | alpha_2   | alpha 2      |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta 1  | beta_1    | beta 1       |
      |        |               | beta 2  | beta_2    | beta 2       |
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | gamma | feature | beta   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | gamma  | local, origin | gamma 1 | gamma_1   | gamma 1      |
      |        |               | gamma 2 | gamma_2   | gamma 2      |
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | delta | feature | gamma  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | delta  | local, origin | delta 1 | delta_1   | delta 1      |
      |        |               | delta 2 | delta_2   | delta 2      |
    And the current branch is "main"
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git checkout alpha       |
      | alpha  | git checkout beta        |
      | beta   | git checkout gamma       |
      | gamma  | git checkout delta       |
      | delta  | git checkout main        |
      | main   | git push --tags          |
    And the initial branches and lineage exist now
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    And the initial branches and lineage exist now
    And the initial commits exist now
