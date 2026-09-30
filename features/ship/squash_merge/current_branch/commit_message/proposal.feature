Feature: use the proposal title, number, and body as the commit message when squash-merging a branch that has a proposal

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | TITLE            | BODY          | URL                      |
      | 1  | feature       | main          | feature proposal | proposal body | https://example.com/pr/1 |
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "feature"

  Scenario: result
    When I run "git-town ship"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                           |
      | feature | git fetch --prune --tags                                          |
      |         | Finding proposal from feature into main ... #1 (feature proposal) |
      |         | git checkout main                                                 |
      | main    | git merge --squash --ff feature                                   |
      |         | git commit -m "feature proposal (#1)                              |
      |         | git push                                                          |
      |         | git push origin :feature                                          |
      |         | git branch -D feature                                             |
    And no lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE               |
      | main   | local, origin | feature proposal (#1) |
    And the initial proposals exist now

  Scenario: an explicitly given commit message still wins
    When I run "git-town ship -m 'custom message'"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit -m "custom message"  |
      |         | git push                        |
      |         | git push origin :feature        |
      |         | git branch -D feature           |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | custom message |
    And the initial proposals exist now

  @skipWindows
  Scenario: the "--enter-message" flag lets the user edit the commit message
    When I run "git-town ship --enter-message" and enter "my message" for the commit message
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                           |
      | feature | git fetch --prune --tags                                          |
      |         | Finding proposal from feature into main ... #1 (feature proposal) |
      |         | git checkout main                                                 |
      | main    | git merge --squash --ff feature                                   |
      |         | git commit                                                        |
      |         | git push                                                          |
      |         | git push origin :feature                                          |
      |         | git branch -D feature                                             |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE    |
      | main   | local, origin | my message |
    And the initial proposals exist now
