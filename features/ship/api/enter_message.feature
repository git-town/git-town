Feature: choosing the commit message when shipping via the forge API

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
      | ID | SOURCE BRANCH | TARGET BRANCH | TITLE            | BODY | URL                      |
      | 1  | feature       | main          | feature proposal |      | https://example.com/pr/1 |
    And Git setting "git-town.ship-strategy" is "api"
    And the current branch is "feature"

  Scenario: default uses the forge's commit message without opening an editor
    When I run "git-town ship"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                           |
      | feature | git fetch --prune --tags                                          |
      |         | Finding proposal from feature into main ... #1 (feature proposal) |
      |         | git checkout main                                                 |
      |         | GitHub API: merging PR #1 ... ok                                  |
      | main    | git push origin :feature                                          |
      |         | git branch -D feature                                             |
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And the initial proposals exist now

  @skipWindows
  Scenario: the "--enter-message" flag lets the user enter the commit message
    When I run "git-town ship --enter-message" and enter "my message" for the commit message
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                           |
      | feature | git fetch --prune --tags                                          |
      |         | Finding proposal from feature into main ... #1 (feature proposal) |
      |         | git checkout main                                                 |
      | main    | git merge --squash --ff feature                                   |
      |         | git commit                                                        |
      |         | git reset --hard HEAD~1                                           |
      |         | GitHub API: merging PR #1 ... ok                                  |
      |         | git push origin :feature                                          |
      |         | git branch -D feature                                             |
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And the initial proposals exist now
