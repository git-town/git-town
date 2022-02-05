Feature: in a local repo

  Background:
    Given my repo does not have a remote origin
    And my repo has the local feature branches "feature" and "other"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
      | other   | local    | other commit   |
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
    And the existing branches are
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
    And my repo now has the commits
      | BRANCH | LOCATION | MESSAGE      |
      | other  | local    | other commit |
    And Git Town now knows this branch hierarchy
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch feature {{ sha 'WIP on feature' }} |
      |         | git checkout feature                          |
      | feature | git reset {{ sha 'feature commit' }}          |
    And I am now on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my initial commits
    And my repo now has its initial branches and branch hierarchy
