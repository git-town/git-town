Feature: git prune-branches: delete branches that were shipped or removed on another machine

  As a developer checking out branches that are also developed on another machine
  I want to remove all branches that have been shipped or deleted on another machine
  So that I keep my local repository free from obsolete branches and remain efficient.

  Rules:
  - branches that have the origin "[<remote name>: gone]" when running "git branch -vv" are obsolete


  Background:
    Given I have feature branches named "active-feature" and "deleted-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                |
      | active-feature  | local and remote | active feature commit  |
      | deleted-feature | local and remote | deleted feature commit |
    And the "deleted-feature" branch gets deleted on the remote
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git prune-branches`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                       |
      | main   | git fetch --prune             |
      |        | git branch -D deleted-feature |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES             |
      | local      | main, active-feature |
      | remote     | main, active-feature |


  Scenario: undo
    When I run `git prune-branches --undo`
    Then it runs the commands
      | BRANCH          | COMMAND                                                        |
      | main            | git branch deleted-feature <%= sha 'deleted feature commit' %> |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                              |
      | local      | main, active-feature, deleted-feature |
      | remote     | main, active-feature                  |
