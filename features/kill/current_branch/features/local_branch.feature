Feature: killing a local branch

  Background:
    Given my repo has a feature branch "feature"
    And my repo has a local feature branch "local-feature"
    And my repo contains the commits
      | BRANCH        | LOCATION      | MESSAGE              |
      | feature       | local, remote | feature commit       |
      | local-feature | local         | local feature commit |
    And I am on the "local-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                              |
      | local-feature | git fetch --prune --tags             |
      |               | git add -A                           |
      |               | git commit -m "WIP on local-feature" |
      |               | git checkout main                    |
      | main          | git branch -D local-feature          |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES      |
      | local, remote | main, feature |
    And my repo now has the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, remote | feature commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH        | COMMAND                                                   |
      | main          | git branch local-feature {{ sha 'WIP on local-feature' }} |
      |               | git checkout local-feature                                |
      | local-feature | git reset {{ sha 'local feature commit' }}                |
    And I am now on the "local-feature" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
    And my repo now has its initial branches and branch hierarchy
