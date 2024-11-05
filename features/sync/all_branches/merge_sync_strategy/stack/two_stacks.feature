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
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | main   | git fetch --prune --tags                |
      |        | git add -A                              |
      |        | git stash                               |
      |        | git rebase origin/main --no-update-refs |
      |        | git checkout first                      |
      | first  | git merge --no-edit --ff main           |
      |        | git merge --no-edit --ff origin/first   |
      |        | git checkout second                     |
      | second | git merge --no-edit --ff first          |
      |        | git merge --no-edit --ff origin/second  |
      |        | git checkout third                      |
      | third  | git merge --no-edit --ff second         |
      |        | git merge --no-edit --ff origin/third   |
      |        | git checkout fourth                     |
      | fourth | git merge --no-edit --ff third          |
      |        | git merge --no-edit --ff origin/fourth  |
      |        | git checkout one                        |
      | one    | git merge --no-edit --ff main           |
      |        | git merge --no-edit --ff origin/one     |
      |        | git checkout two                        |
      | two    | git merge --no-edit --ff one            |
      |        | git merge --no-edit --ff origin/two     |
      |        | git checkout three                      |
      | three  | git merge --no-edit --ff two            |
      |        | git merge --no-edit --ff origin/three   |
      |        | git checkout four                       |
      | four   | git merge --no-edit --ff three          |
      |        | git merge --no-edit --ff origin/four    |
      |        | git checkout main                       |
      | main   | git push --tags                         |
      |        | git stash pop                           |
    And the current branch is still "main"
    And the uncommitted file still exists
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | main   | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "main"
    And the uncommitted file still exists
    And the initial commits exist now
    And the initial branches and lineage exist now
