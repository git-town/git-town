Feature: sync a stack making independent changes

  Background:
    Given a Git repo clone
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
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | main   | git fetch --prune --tags              |
      |        | git add -A                            |
      |        | git stash                             |
      |        | git rebase origin/main                |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
      |        | git merge --no-edit --ff main         |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff origin/beta  |
      |        | git merge --no-edit --ff alpha        |
      |        | git checkout gamma                    |
      | gamma  | git merge --no-edit --ff origin/gamma |
      |        | git merge --no-edit --ff beta         |
      |        | git checkout delta                    |
      | delta  | git merge --no-edit --ff origin/delta |
      |        | git merge --no-edit --ff gamma        |
      |        | git checkout main                     |
      | main   | git push --tags                       |
      |        | git stash pop                         |
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
