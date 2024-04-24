Feature: on a feature branch in a repository with a submodule that has uncommitted changes

  Background:
    Given my repo has a Git submodule
    And the current branch is a feature branch "feature"
    And an uncommitted file with name "submodule/file" and content "a change in the submodule"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit --ff main      |
    And the current branch is still "feature"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE         |
      | main    | local, origin | added submodule |
      | feature | local, origin | added submodule |
    And the initial branches and lineage exist
