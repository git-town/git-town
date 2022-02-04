Feature: auto-push new branches

  Background:
    Given the new-branch-push-flag configuration is true
    And my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, remote | feature_commit |
    And I am on the "feature" branch
    When I run "git-town prepend parent"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                   |
      | feature | git fetch --prune --tags  |
      |         | git checkout main         |
      | main    | git rebase origin/main    |
      |         | git branch parent main    |
      |         | git checkout parent       |
      | parent  | git push -u origin parent |
    And I am now on the "parent" branch
    And my repo now has the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, remote | feature_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | parent |
      | parent  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                 |
      | parent | git push origin :parent |
      |        | git checkout main       |
      | main   | git branch -d parent    |
      |        | git checkout feature    |
    And I am now on the "feature" branch
    And my repo is left with my original commits
    And Git Town now has the original branch hierarchy
