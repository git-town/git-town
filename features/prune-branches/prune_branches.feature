Feature: delete branches that were shipped or removed on another machine

  Background:
    Given my repo has the feature branches "active-feature" and "finished-feature"
    And the following commits exist in my repo
      | BRANCH           | LOCATION      | MESSAGE                 |
      | active-feature   | local, remote | active-feature commit   |
      | finished-feature | local, remote | finished-feature commit |
    And the "finished-feature" branch gets deleted on the remote
    And I am on the "finished-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town prune-branches"

  Scenario: result
    Then it runs the commands
      | BRANCH           | COMMAND                        |
      | finished-feature | git fetch --prune --tags       |
      |                  | git checkout main              |
      | main             | git branch -D finished-feature |
    And I am now on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES             |
      | local      | main, active-feature |
      | remote     | main, active-feature |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT |
      | active-feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                         |
      | main   | git branch finished-feature {{ sha 'finished-feature commit' }} |
      |        | git checkout finished-feature                                   |
    And I am now on the "finished-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                               |
      | local      | main, active-feature, finished-feature |
      | remote     | main, active-feature                   |
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT |
      | active-feature   | main   |
      | finished-feature | main   |
