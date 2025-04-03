Feature: merging a branch with uncommitted changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta commit | beta-file | beta content |
    And the current branch is "beta"
    And an uncommitted file
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | beta   | git fetch --prune --tags    |
      |        | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git branch -D alpha         |
      |        | git push origin :alpha      |
      |        | git stash pop               |
      |        | git restore --staged .      |
    And the current branch is still "beta"
    And this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | beta   | local, origin | alpha commit | alpha-file | alpha content |
      |        |               | beta commit  | beta-file  | beta content  |
    And these committed files exist now
      | BRANCH | NAME       | CONTENT       |
      | beta   | alpha-file | alpha content |
      |        | beta-file  | beta content  |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git add -A                                           |
      |        | git stash -m "Git Town WIP"                          |
      |        | git branch alpha {{ sha-before-run 'alpha commit' }} |
      |        | git push -u origin alpha                             |
      |        | git stash pop                                        |
      |        | git restore --staged .                               |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
    And the uncommitted file still exists
