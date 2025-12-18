Feature: ignore uncommitted changes using the config file

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "feature"
    And the configuration file:
      """
      [ship]
      ignore-uncommitted = true
      """
    And an uncommitted file
    When I run "git-town ship -m shipped"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit -m shipped           |
      |         | git push                        |
      |         | git push origin :feature        |
      |         | git branch -D feature           |
    And the current branch is now "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE |
      | main   | local, origin | shipped |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git add -A                                    |
      |         | git stash -m "Git Town WIP"                   |
      |         | git revert {{ sha 'shipped' }}                |
      |         | git push                                      |
      |         | git branch feature {{ sha 'feature commit' }} |
      |         | git push -u origin feature                    |
      |         | git checkout feature                          |
      | feature | git stash pop                                 |
      |         | git restore --staged .                        |
    And the initial lineage exists now
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE          |
      | main    | local, origin | shipped          |
      |         |               | Revert "shipped" |
      | feature | local, origin | feature commit   |
    And the uncommitted file still exists
