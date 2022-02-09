Feature: in a local repo

  Background:
    Given my repo does not have an origin
    And the local feature branches "feature" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
      | other   | local    | other commit   |
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                        |
      | feature | git add -A                     |
      |         | git commit -m "WIP on feature" |
      |         | git checkout main              |
      | main    | git branch -D feature          |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
    And now these commits exist
      | BRANCH | LOCATION | MESSAGE      |
      | other  | local    | other commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch feature {{ sha 'WIP on feature' }} |
      |         | git checkout feature                          |
      | feature | git reset {{ sha 'feature commit' }}          |
    And the current branch is now "feature"
    And my workspace still contains my uncommitted file
    And now the initial commits exist
    And the initial branches and hierarchy exist
