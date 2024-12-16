@skipWindows
Feature: ship to a custom dev remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And I rename the "origin" remote to "fork"
    And Git Town setting "dev-remote" is "fork"
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship" and enter "feature done" for the commit message

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit                      |
      |         | git push                        |
      |         | git push fork :feature          |
      |         | git branch -D feature           |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY  | BRANCHES |
      | local, fork | main     |
    And these commits exist now
      | BRANCH | LOCATION    | MESSAGE      |
      | main   | local, fork | feature done |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git revert {{ sha 'feature done' }}           |
      |        | git push                                      |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git push -u fork feature                      |
      |        | git checkout feature                          |
    And the current branch is now "feature"
    And these commits exist now
      | BRANCH  | LOCATION    | MESSAGE               |
      | main    | local, fork | feature done          |
      |         |             | Revert "feature done" |
      | feature | local, fork | feature commit        |
    And the initial branches and lineage exist now
