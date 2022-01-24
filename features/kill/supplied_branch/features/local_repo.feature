Feature: local repository

  Background:
    Given my repo does not have a remote origin
    And my repo has the local feature branches "dead-feature" and "other-feature"
    And the following commits exist in my repo
      | BRANCH        | LOCATION | MESSAGE              | FILE NAME            |
      | main          | local    | main commit          | conflicting_file     |
      | dead-feature  | local    | dead feature commit  | current_feature_file |
      | other-feature | local    | other feature commit | other_feature_file   |
    And I am on the "dead-feature" branch
    And my workspace has an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town kill other-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                     |
      | dead-feature | git add -A                  |
      |              | git stash                   |
      |              | git branch -D other-feature |
      |              | git stash pop               |
    And I am still on the "dead-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES           |
      | local      | main, dead-feature |
    And my repo now has the following commits
      | BRANCH       | LOCATION | MESSAGE             | FILE NAME            |
      | main         | local    | main commit         | conflicting_file     |
      | dead-feature | local    | dead feature commit | current_feature_file |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH       | COMMAND                                                   |
      | dead-feature | git add -A                                                |
      |              | git stash                                                 |
      |              | git branch other-feature {{ sha 'other feature commit' }} |
      |              | git stash pop                                             |
    And I am still on the "dead-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                          |
      | local      | main, dead-feature, other-feature |
    And my repo is left with my original commits
