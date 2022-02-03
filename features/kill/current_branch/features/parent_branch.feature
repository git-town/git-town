Feature: killing a branch within a branch chain

  Background:
    Given my repo has a feature branch "feature-1"
    And my repo has a feature branch "feature-2" as a child of "feature-1"
    And my repo has a feature branch "feature-3" as a child of "feature-2"
    And my repo contains the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature-1 | local, remote | feature 1 commit |
      | feature-2 | local, remote | feature 2 commit |
      | feature-3 | local, remote | feature 3 commit |
    And I am on the "feature-2" branch
    And my workspace has an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                          |
      | feature-2 | git fetch --prune --tags         |
      |           | git push origin :feature-2       |
      |           | git add -A                       |
      |           | git commit -m "WIP on feature-2" |
      |           | git checkout feature-1           |
      | feature-1 | git branch -D feature-2          |
    And I am now on the "feature-1" branch
    And my repo doesn't have any uncommitted files
    And the existing branches are
      | REPOSITORY    | BRANCHES                   |
      | local, remote | main, feature-1, feature-3 |
    And my repo now has the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature-1 | local, remote | feature 1 commit |
      | feature-3 | local, remote | feature 3 commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-3 | feature-1 |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH    | COMMAND                                           |
      | feature-1 | git branch feature-2 {{ sha 'WIP on feature-2' }} |
      |           | git checkout feature-2                            |
      | feature-2 | git reset {{ sha 'feature 2 commit' }}            |
      |           | git push -u origin feature-2                      |
    And I am now on the "feature-2" branch
    And my workspace has the uncommitted file again
    And my repo is left with my original commits
    And my repo now has its original branches and branch hierarchy
