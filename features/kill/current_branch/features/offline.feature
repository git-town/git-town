Feature: offline mode

  Background:
    Given Git Town is in offline mode
    And my repo has the feature branches "feature" and "other-feature"
    And my repo contains the commits
      | BRANCH        | LOCATION      | MESSAGE              |
      | feature       | local, remote | feature commit       |
      | other-feature | local, remote | other feature commit |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                        |
      | feature | git add -A                     |
      |         | git commit -m "WIP on feature" |
      |         | git checkout main              |
      | main    | git branch -D feature          |
    And I am now on the "main" branch
    And my repo doesn't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES                     |
      | local      | main, other-feature          |
      | remote     | main, feature, other-feature |
    And my repo now has the following commits
      | BRANCH        | LOCATION      | MESSAGE              |
      | feature       | remote        | feature commit       |
      | other-feature | local, remote | other feature commit |
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
    And the existing branches are
      | REPOSITORY    | BRANCHES                     |
      | local, remote | main, feature, other-feature |
    And my repo is left with my original commits
    And Git Town now has the original branch hierarchy
