Feature: sync a workspace with two independent stacks

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | one  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE |
      | one    | local, origin | one     |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | two  | feature | one    | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE |
      | two    | local, origin | two     |
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | three | feature | two    | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE |
      | three  | local, origin | three   |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | four | feature | three  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE |
      | four   | local, origin | four    |
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | first | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE |
      | first  | local, origin | first 1 |
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | second | feature | first  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  |
      | second | local, origin | second 1 |
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | third | feature | second | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE |
      | third  | local, origin | third 1 |
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | fourth | feature | third  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  |
      | fourth | local, origin | fourth 1 |
    And the current branch is "main"
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git checkout first       |
      | first  | git checkout second      |
      | second | git checkout third       |
      | third  | git checkout fourth      |
      | fourth | git checkout one         |
      | one    | git checkout two         |
      | two    | git checkout three       |
      | three  | git checkout four        |
      | four   | git checkout main        |
      | main   | git push --tags          |
    And the initial branches and lineage exist now
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    And the initial branches and lineage exist now
    And the initial commits exist now
