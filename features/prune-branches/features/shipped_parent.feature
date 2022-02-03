Feature: a parent branch of a local branch was shipped

  Background:
    Given my repo has a feature branch "feature"
    And my repo has a feature branch "feature-child" as a child of "feature"
    And my repo contains the commits
      | BRANCH        | LOCATION      | MESSAGE              |
      | feature       | local, remote | feature commit       |
      | feature-child | local, remote | feature-child commit |
    And the "feature" branch gets deleted on the remote
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town prune-branches"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git branch -D feature    |
    And I am now on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY    | BRANCHES            |
      | local, remote | main, feature-child |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | feature-child | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'feature commit' }} |
    And I am now on the "main" branch
    And my workspace still contains my uncommitted file
    And my repo now has its original branches and branch hierarchy
