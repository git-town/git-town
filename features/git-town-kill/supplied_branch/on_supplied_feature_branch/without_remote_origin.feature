Feature: git town-kill: killing the given feature branch when on it (without remote repo)

  (see ../with_tracking_branch/with_open_changes.feature)


  Background:
    Given my repo does not have a remote origin
    And my repository has the local feature branches "current-feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION | MESSAGE                |
      | current-feature | local    | current feature commit |
      | other-feature   | local    | other feature commit   |
    And I am on the "current-feature" branch
    And my workspace has an uncommitted file
    When I run `git-town kill current-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                                |
      | current-feature | git add -A                             |
      |                 | git commit -m "WIP on current-feature" |
      |                 | git checkout main                      |
      | main            | git branch -D current-feature          |
    And I end up on the "main" branch
    And I don't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
    And now my repository has the following commits
      | BRANCH        | LOCATION | MESSAGE              |
      | other-feature | local    | other feature commit |


  Scenario: undoing the kill
    When I run `git-town undo`
    Then it runs the commands
      | BRANCH          | COMMAND                                                        |
      | main            | git branch current-feature <%= sha 'WIP on current-feature' %> |
      |                 | git checkout current-feature                                   |
      | current-feature | git reset <%= sha 'current feature commit' %>                  |
    And I end up on the "current-feature" branch
    And my workspace has the uncommitted file again
    And the existing branches are
      | REPOSITORY | BRANCHES                             |
      | local      | main, current-feature, other-feature |
    And my repository is left with my original commits
