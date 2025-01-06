Feature: sync all feature branches in the presence of uncommitted changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git add -A                              |
      |         | git stash -m "Git Town WIP"             |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git push --tags                         |
      |         | git stash pop                           |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git add -A                  |
      |         | git stash -m "Git Town WIP" |
      |         | git stash pop               |
    And the current branch is still "feature"
    And the initial commits exist now
    And the uncommitted file still exists
