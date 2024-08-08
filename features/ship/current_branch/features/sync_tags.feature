@smoke
Feature: don't sync tags while shipping

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "feature"
    And Git Town setting "sync-tags" is "false"
    When I run "git-town ship -m done"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --no-tags     |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit -m done              |
      |         | git push                        |
      |         | git push origin :feature        |
      |         | git branch -D feature           |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git revert {{ sha 'done' }}                   |
      |        | git push                                      |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git push -u origin feature                    |
      |        | git checkout feature                          |
    And the initial tags exist now
