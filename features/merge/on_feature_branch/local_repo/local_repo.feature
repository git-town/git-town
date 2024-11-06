Feature: merging in a local repo

  Background:
    Given a local Git repo
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | alpha | feature | main   | local     |
      | beta  | feature | alpha  | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | alpha  | local    | alpha commit |
      | beta   | local    | beta commit  |
    And the current branch is "beta"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                        |
      | beta   | git merge --no-edit --ff alpha |
      |        | git branch -D alpha            |
    And the current branch is still "beta"
    And this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE                        |
      | beta   | local    | beta commit                    |
      |        |          | alpha commit                   |
      |        |          | Merge branch 'alpha' into beta |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git reset --hard {{ sha-before-run 'beta commit' }}  |
      |        | git branch alpha {{ sha-before-run 'alpha commit' }} |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
