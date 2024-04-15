Feature: sync a stack making independent changes

  Background:
    Given feature branch "alpha" with these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | alpha 1 | alpha_1   | alpha 1      |
      |               | alpha 2 | alpha_2   | alpha 2      |
    And feature branch "beta" as a child of "alpha" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | beta 1  | beta_1    | beta 1       |
      |               | beta 2  | beta_2    | beta 2       |
    And feature branch "gamma" as a child of "beta" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | gamma 1 | gamma_1   | gamma 1      |
      |               | gamma 2 | gamma_2   | gamma 2      |
    And feature branch "delta" as a child of "gamma" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | delta 1 | delta_1   | delta 1      |
      |               | delta 2 | delta_2   | delta 2      |
    And the current branch is "delta"
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                          |
      | delta  | git fetch --prune --tags         |
      |        | git add -A                       |
      |        | git stash                        |
      |        | git checkout main                |
      | main   | git rebase origin/main           |
      |        | git checkout alpha               |
      | alpha  | git merge --no-edit origin/alpha |
      |        | git merge --no-edit main         |
      |        | git checkout beta                |
      | beta   | git merge --no-edit origin/beta  |
      |        | git merge --no-edit alpha        |
      |        | git checkout gamma               |
      | gamma  | git merge --no-edit origin/gamma |
      |        | git merge --no-edit beta         |
      |        | git checkout delta               |
      | delta  | git merge --no-edit origin/delta |
      |        | git merge --no-edit gamma        |
      |        | git push --tags                  |
      |        | git stash pop                    |
    And the current branch is still "delta"
    And the initial commits exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | delta  | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "delta"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
