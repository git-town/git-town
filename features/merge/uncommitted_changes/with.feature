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
      | BRANCH | COMMAND                                         |
      | beta   | git fetch --prune --tags                        |
      |        | git add -A                                      |
      |        | git stash -m "Git Town WIP"                     |
      |        | git checkout alpha                              |
      | alpha  | git reset --hard {{ sha 'beta commit' }}        |
      |        | git push origin :beta                           |
      |        | git branch -D beta                              |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
      |        | git restore --staged .                          |
    And this lineage exists now
      """
      main
        alpha
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
      |        |               | beta commit  | beta-file  | beta content  |
    And these committed files exist now
      | BRANCH | NAME       | CONTENT       |
      | alpha  | alpha-file | alpha content |
      |        | beta-file  | beta content  |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git add -A                                      |
      |        | git stash -m "Git Town WIP"                     |
      |        | git reset --hard {{ sha 'alpha commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git branch beta {{ sha 'beta commit' }}         |
      |        | git push -u origin beta                         |
      |        | git checkout beta                               |
      | beta   | git stash pop                                   |
      |        | git restore --staged .                          |
    And the initial lineage exists now
    And the uncommitted file still exists
    And the initial commits exist now
