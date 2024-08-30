# TODO: delete this
@skipWindows
Feature: ship a branch using the "compress" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE          |
      | feature | local, origin | feature commit 1 |
      | feature | local, origin | feature commit 2 |
    And Git Town setting "sync-feature-strategy" is "compress"
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship" and enter "feature done" for the commit message

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit                      |
      |         | git push                        |
      |         | git push origin :feature        |
      |         | git branch -D feature           |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | feature done |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | main   | git revert {{ sha 'feature done' }}             |
      |        | git push                                        |
      |        | git branch feature {{ sha 'feature commit 2' }} |
      |        | git push -u origin feature                      |
      |        | git checkout feature                            |
    And the current branch is now "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | feature done          |
      |         |               | Revert "feature done" |
      | feature | local, origin | feature commit 1      |
      |         |               | feature commit 2      |
    And the initial branches and lineage exist
