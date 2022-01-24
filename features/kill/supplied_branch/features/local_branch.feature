Feature: local branch

  Background:
    Given my repo does not have a remote origin
    And my repo has the local feature branches "dead-feature" and "other-feature"
    And the following commits exist in my repo
      | BRANCH        | LOCATION | MESSAGE              |
      | dead-feature  | local    | dead feature commit  |
      | other-feature | local    | other feature commit |
    And I am on the "dead-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town kill dead-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                             |
      | dead-feature | git add -A                          |
      |              | git commit -m "WIP on dead-feature" |
      |              | git checkout main                   |
      | main         | git branch -D dead-feature          |
    And I am now on the "main" branch
    And my repo doesn't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
    And my repo now has the following commits
      | BRANCH        | LOCATION | MESSAGE              |
      | other-feature | local    | other feature commit |

  Scenario: undoing the kill
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH       | COMMAND                                                 |
      | main         | git branch dead-feature {{ sha 'WIP on dead-feature' }} |
      |              | git checkout dead-feature                               |
      | dead-feature | git reset {{ sha 'dead feature commit' }}               |
    And I am now on the "dead-feature" branch
    And my workspace has the uncommitted file again
    And the existing branches are
      | REPOSITORY | BRANCHES                          |
      | local      | main, dead-feature, other-feature |
    And my repo is left with my original commits
