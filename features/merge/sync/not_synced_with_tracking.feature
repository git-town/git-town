Feature: merging when the branch is not in sync with its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE            |
      | beta   | local    | local beta commit  |
      | beta   | origin   | remote beta commit |
    And the current branch is "beta"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                   |
      | beta   | git merge --abort                                         |
      |        | git reset --hard {{ sha-before-run 'local beta commit' }} |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file" with "resolved beta content"
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH | COMMAND                |
      | beta   | git commit --no-edit   |
      |        | git push               |
      |        | git branch -D alpha    |
      |        | git push origin :alpha |
    And the current branch is still "beta"
    And this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                              | FILE NAME        | FILE CONTENT          |
      | beta   | local, origin | local beta commit                                    | conflicting_file | local beta content    |
      |        |               | alpha commit                                         | alpha_file       | alpha content         |
      |        |               | Merge branch 'alpha' into beta                       |                  |                       |
      |        |               | remote beta commit                                   | conflicting_file | remote beta content   |
      |        |               | Merge remote-tracking branch 'origin/beta' into beta | conflicting_file | resolved beta content |
