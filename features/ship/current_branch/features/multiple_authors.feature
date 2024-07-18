@messyoutput
@skipWindows
Feature: ship a coworker's feature branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE            | AUTHOR                            |
      | feature | local    | developer commit 1 | developer <developer@example.com> |
      |         |          | developer commit 2 | developer <developer@example.com> |
      |         |          | coworker commit    | coworker <coworker@example.com>   |

  Scenario: choose myself as the author
    When I run "git-town ship -m 'feature done'" and enter into the dialog:
      | DIALOG                              | KEYS  |
      | choose author for the squash commit | enter |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                            |
      | main   | local, origin | feature done | developer <developer@example.com> |
    And no lineage exists now

  Scenario: choose a coworker as the author
    When I run "git-town ship -m 'feature done'" and enter into the dialog:
      | DIALOG                              | KEYS       |
      | choose author for the squash commit | down enter |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                          |
      | main   | local, origin | feature done | coworker <coworker@example.com> |
    And no lineage exists now

  Scenario: undo
    Given I ran "git-town ship -m 'feature done'" and enter into the dialog:
      | DIALOG                              | KEYS  |
      | choose author for the squash commit | enter |
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                       |
      | main   | git revert {{ sha 'feature done' }}                           |
      |        | git push                                                      |
      |        | git push origin {{ sha 'initial commit' }}:refs/heads/feature |
      |        | git branch feature {{ sha 'coworker commit' }}                |
      |        | git checkout feature                                          |
    And the current branch is now "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | feature done          |
      |         |               | Revert "feature done" |
      | feature | local         | developer commit 1    |
      |         |               | developer commit 2    |
      |         |               | coworker commit       |
    And the initial branches and lineage exist
