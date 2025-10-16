Feature: sync all feature branches in the presence of uncommitted changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
      | feature-2 | feature | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature-1 | local, origin | feature 1 commit |
      | feature-2 | local, origin | feature 2 commit |
    And the current branch is "feature-1"
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                     |
      | feature-1 | git fetch --prune --tags    |
      |           | git add -A                  |
      |           | git stash -m "Git Town WIP" |
      |           | git checkout feature-2      |
      | feature-2 | git checkout feature-1      |
      | feature-1 | git push --tags             |
      |           | git stash pop               |
      |           | git restore --staged .      |
    And the uncommitted file still exists
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                     |
      | feature-1 | git add -A                  |
      |           | git stash -m "Git Town WIP" |
      |           | git stash pop               |
      |           | git restore --staged .      |
    And the uncommitted file still exists
    And the initial commits exist now
