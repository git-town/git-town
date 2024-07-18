Feature: ship a coworker's feature branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE         | AUTHOR                          |
      | feature | local, origin | coworker commit | coworker <coworker@example.com> |

  Scenario: result (commit message via CLI)
    When I run "git-town ship -m 'feature done'"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                 |
      | feature | git fetch --prune --tags                                                |
      |         | git checkout main                                                       |
      | main    | git merge --squash --ff feature                                         |
      |         | git commit -m "feature done" --author "coworker <coworker@example.com>" |
      |         | git push                                                                |
      |         | git push origin :feature                                                |
      |         | git branch -D feature                                                   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                          |
      | main   | local, origin | feature done | coworker <coworker@example.com> |
    And no lineage exists now

  @skipWindows
  Scenario: result (commit message via editor)
    When I run "git-town ship" and enter "feature done" for the commit message
    Then it runs the commands
      | BRANCH  | COMMAND                                               |
      | feature | git fetch --prune --tags                              |
      |         | git checkout main                                     |
      | main    | git merge --squash --ff feature                       |
      |         | git commit --author "coworker <coworker@example.com>" |
      |         | git push                                              |
      |         | git push origin :feature                              |
      |         | git branch -D feature                                 |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                          |
      | main   | local, origin | feature done | coworker <coworker@example.com> |
    And no lineage exists now

  Scenario: undo
    Given I ran "git-town ship -m 'feature done'"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                        |
      | main   | git revert {{ sha 'feature done' }}            |
      |        | git push                                       |
      |        | git branch feature {{ sha 'coworker commit' }} |
      |        | git push -u origin feature                     |
      |        | git checkout feature                           |
    And the current branch is now "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | feature done          |
      |         |               | Revert "feature done" |
      | feature | local, origin | coworker commit       |
    And the initial branches and lineage exist
