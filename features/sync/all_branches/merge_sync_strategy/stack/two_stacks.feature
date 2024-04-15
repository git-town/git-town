Feature: sync a stack making independent changes

  Background:
    Given feature branch "alpha" with these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | alpha 1 | alpha_1   | alpha 1      |
    And feature branch "beta" as a child of "alpha" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | beta 1  | beta_1    | beta 1       |
    And feature branch "gamma" as a child of "beta" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | gamma 1 | gamma_1   | gamma 1      |
    And feature branch "delta" as a child of "gamma" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | delta 1 | delta_1   | delta 1      |
    And feature branch "first" with these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | first 1 | first_1   | first 1      |
    And feature branch "second" as a child of "first" has these commits
      | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | local, origin | second 1 | second_1  | second 1     |
    And feature branch "third" as a child of "second" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | third 1 | third_1   | third 1      |
    And feature branch "fourth" as a child of "third" has these commits
      | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | local, origin | fourth 1 | fourth_1  | fourth 1     |
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                           |
      | main   | git fetch --prune --tags          |
      |        | git add -A                        |
      |        | git stash                         |
      |        | git rebase origin/main            |
      |        | git checkout alpha                |
      | alpha  | git merge --no-edit origin/alpha  |
      |        | git merge --no-edit main          |
      |        | git checkout beta                 |
      | beta   | git merge --no-edit origin/beta   |
      |        | git merge --no-edit alpha         |
      |        | git checkout delta                |
      | delta  | git merge --no-edit origin/delta  |
      |        | git merge --no-edit gamma         |
      |        | git checkout first                |
      | first  | git merge --no-edit origin/first  |
      |        | git merge --no-edit main          |
      |        | git checkout fourth               |
      | fourth | git merge --no-edit origin/fourth |
      |        | git merge --no-edit third         |
      |        | git checkout gamma                |
      | gamma  | git merge --no-edit origin/gamma  |
      |        | git merge --no-edit beta          |
      |        | git checkout second               |
      | second | git merge --no-edit origin/second |
      |        | git merge --no-edit first         |
      |        | git checkout third                |
      | third  | git merge --no-edit origin/third  |
      |        | git merge --no-edit second        |
      |        | git checkout main                 |
      | main   | git push --tags                   |
      |        | git stash pop                     |
    And the current branch is still "main"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | main   | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "main"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
