Feature: merging a branch when everything is conflicting

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE             | FILE NAME        | FILE CONTENT         |
      | alpha  | local    | local alpha commit  | conflicting_file | local alpha content  |
      | alpha  | origin   | remote alpha commit | conflicting_file | remote alpha content |
      | beta   | local    | local beta commit   | conflicting_file | local beta content   |
      | beta   | origin   | remote beta commit  | conflicting_file | remote beta content  |
    And the current branch is "beta"
    When I run "git-town merge"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | beta   | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | alpha  | git merge --abort |
      |        | git checkout beta |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now

  Scenario: resolve all conflicts and undo
    When I resolve the conflict in "conflicting_file" with "resolved between local and tracking alpha"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND                        |
      | alpha  | git commit --no-edit           |
      |        | git checkout beta              |
      | beta   | git merge --no-edit --ff alpha |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When I resolve the conflict in "conflicting_file" with "resolved between local beta and local alpha"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND                              |
      | beta   | git commit --no-edit                 |
      |        | git merge --no-edit --ff origin/beta |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When I resolve the conflict in "conflicting_file" with "resolved between local and tracking beta"
    And I run "git-town continue" and close the editor
    Then it runs the commands
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
      | BRANCH | LOCATION      | MESSAGE                                                | FILE NAME        | FILE CONTENT                                |
      | beta   | local, origin | local beta commit                                      | conflicting_file | local beta content                          |
      |        |               | local alpha commit                                     | conflicting_file | local alpha content                         |
      |        |               | remote alpha commit                                    | conflicting_file | remote alpha content                        |
      |        |               | Merge remote-tracking branch 'origin/alpha' into alpha | conflicting_file | resolved between local and tracking alpha   |
      |        |               | Merge branch 'alpha' into beta                         | conflicting_file | resolved between local beta and local alpha |
      |        |               | remote beta commit                                     | conflicting_file | remote beta content                         |
      |        |               | Merge remote-tracking branch 'origin/beta' into beta   | conflicting_file | resolved between local and tracking beta    |
    And these committed files exist now
      | BRANCH | NAME             | CONTENT                                  |
      | beta   | conflicting_file | resolved between local and tracking beta |
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                          |
      | beta   | git reset --hard {{ sha-before-run 'local beta commit' }}                        |
      |        | git push --force-with-lease origin {{ sha-in-origin 'remote beta commit' }}:beta |
      |        | git push origin {{ sha-in-origin 'remote alpha commit' }}:refs/heads/alpha       |
      |        | git branch alpha {{ sha-before-run 'local alpha commit' }}                       |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
