Feature: sync a workspace with two independent stacks

  Background:
    Given feature branch "one" with these commits
      | LOCATION      | MESSAGE |
      | local, origin | one     |
    And feature branch "two" as a child of "one" has these commits
      | LOCATION      | MESSAGE |
      | local, origin | two     |
    And feature branch "three" as a child of "two" has these commits
      | LOCATION      | MESSAGE |
      | local, origin | three   |
    And feature branch "four" as a child of "three" has these commits
      | LOCATION      | MESSAGE |
      | local, origin | four    |
    And feature branch "first" with these commits
      | LOCATION      | MESSAGE |
      | local, origin | first 1 |
    And feature branch "second" as a child of "first" has these commits
      | LOCATION      | MESSAGE  |
      | local, origin | second 1 |
    And feature branch "third" as a child of "second" has these commits
      | LOCATION      | MESSAGE |
      | local, origin | third 1 |
    And feature branch "fourth" as a child of "third" has these commits
      | LOCATION      | MESSAGE  |
      | local, origin | fourth 1 |
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                |
      | main   | git fetch --prune --tags               |
      |        | git add -A                             |
      |        | git stash                              |
      |        | git rebase origin/main                 |
      |        | git checkout first                     |
      | first  | git merge --no-edit --ff origin/first  |
      |        | git merge --no-edit --ff main          |
      |        | git checkout second                    |
      | second | git merge --no-edit --ff origin/second |
      |        | git merge --no-edit --ff first         |
      |        | git checkout third                     |
      | third  | git merge --no-edit --ff origin/third  |
      |        | git merge --no-edit second             |
      |        | git checkout fourth                    |
      | fourth | git merge --no-edit --ff origin/fourth |
      |        | git merge --no-edit third              |
      |        | git checkout one                       |
      | one    | git merge --no-edit --ff origin/one    |
      |        | git merge --no-edit --ff main          |
      |        | git checkout two                       |
      | two    | git merge --no-edit --ff origin/two    |
      |        | git merge --no-edit --ff one           |
      |        | git checkout three                     |
      | three  | git merge --no-edit --ff origin/three  |
      |        | git merge --no-edit two                |
      |        | git checkout four                      |
      | four   | git merge --no-edit --ff origin/four   |
      |        | git merge --no-edit three              |
      |        | git checkout main                      |
      | main   | git push --tags                        |
      |        | git stash pop                          |
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
