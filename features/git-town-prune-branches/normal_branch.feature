Feature: git town-prune-branches: delete branches that were shipped or removed on another machine

  As a developer checking out branches that are also developed on another machine
  I want to remove all branches that have been shipped or deleted on another machine
  So that I keep my local repository free from obsolete branches and remain efficient.

  Rules:
  - branches with a deleted tracking branch are removed
  - "git branch -vv" shows these branches with the remote branch name as "[origin/<branch name>: gone]"


  Background:
    Given my repository has the feature branches "active-feature" and "deleted-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                |
      | active-feature  | local and remote | active-feature commit  |
      | deleted-feature | local and remote | deleted-feature commit |
    And the "deleted-feature" branch gets deleted on the remote
    And I am on the "deleted-feature" branch
    And my workspace has an uncommitted file
    When I run `git-town prune-branches`


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                       |
      | deleted-feature | git fetch --prune             |
      |                 | git checkout main             |
      | main            | git branch -D deleted-feature |
    And I end up on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES             |
      | local      | main, active-feature |
      | remote     | main, active-feature |


  Scenario: undo
    When I run `git-town undo`
    Then it runs the commands
      | BRANCH | COMMAND                                                        |
      | main   | git branch deleted-feature <%= sha 'deleted-feature commit' %> |
      |        | git checkout deleted-feature                                   |
    And I end up on the "deleted-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                              |
      | local      | main, active-feature, deleted-feature |
      | remote     | main, active-feature                  |
