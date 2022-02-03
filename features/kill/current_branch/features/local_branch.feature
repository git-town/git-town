Feature: killing a local branch

  Background:
    Given my repo has a feature branch "other-feature"
    And my repo has a local feature branch "current-feature"
    And my repo contains the commits
      | BRANCH          | LOCATION      | MESSAGE                |
      | current-feature | local         | current feature commit |
      | other-feature   | local, remote | other feature commit   |
    And I am on the "current-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                                |
      | current-feature | git fetch --prune --tags               |
      |                 | git add -A                             |
      |                 | git commit -m "WIP on current-feature" |
      |                 | git checkout main                      |
      | main            | git branch -D current-feature          |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES            |
      | local, remote | main, other-feature |
    And my repo now has the following commits
      | BRANCH        | LOCATION      | MESSAGE              |
      | other-feature | local, remote | other feature commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | other-feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH          | COMMAND                                                       |
      | main            | git branch current-feature {{ sha 'WIP on current-feature' }} |
      |                 | git checkout current-feature                                  |
      | current-feature | git reset {{ sha 'current feature commit' }}                  |
    And I am now on the "current-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the initial branches
    And my repo is left with my original commits
    And Git Town now has the original branch hierarchy
