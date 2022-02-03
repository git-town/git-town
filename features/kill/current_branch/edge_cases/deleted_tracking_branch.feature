Feature: the branch to kill has a deleted tracking branch

  Background:
    Given my repo has the feature branches "feature" and "other-feature"
    And my repo contains the commits
      | BRANCH        | LOCATION      | MESSAGE              |
      | feature       | local, remote | feature commit       |
      | other-feature | local, remote | other feature commit |
    And the "feature" branch gets deleted on the remote
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                        |
      | feature | git fetch --prune --tags       |
      |         | git add -A                     |
      |         | git commit -m "WIP on feature" |
      |         | git checkout main              |
      | main    | git branch -D feature          |
    And I am now on the "main" branch
    And my repo doesn't have any uncommitted files
    And the existing branches are
      | REPOSITORY    | BRANCHES            |
      | local, remote | main, other-feature |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | other-feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch feature {{ sha 'WIP on feature' }} |
      |         | git checkout feature                          |
      | feature | git reset {{ sha 'feature commit' }}          |
    And I am now on the "feature" branch
    And my workspace has the uncommitted file again
    And my repo now has the initial branches
    And Git Town now has the original branch hierarchy
