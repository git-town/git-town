Feature: local repository

  Background:
    Given my repo does not have a remote origin
    And my repo has the local feature branches "good-feature" and "other-feature"
    And my repo contains the commits
      | BRANCH        | LOCATION | MESSAGE              | FILE NAME        |
      | main          | local    | main commit          | conflicting_file |
      | good-feature  | local    | good feature commit  | file             |
      | other-feature | local    | other feature commit | file             |
    And I am on the "good-feature" branch
    And my workspace has an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town kill other-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                     |
      | good-feature | git add -A                  |
      |              | git stash                   |
      |              | git branch -D other-feature |
      |              | git stash pop               |
    And I am still on the "good-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES           |
      | local      | main, good-feature |
    And my repo now has the commits
      | BRANCH       | LOCATION | MESSAGE             |
      | main         | local    | main commit         |
      | good-feature | local    | good feature commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH       | PARENT |
      | good-feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH       | COMMAND                                                   |
      | good-feature | git add -A                                                |
      |              | git stash                                                 |
      |              | git branch other-feature {{ sha 'other feature commit' }} |
      |              | git stash pop                                             |
    And I am still on the "good-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the initial branches
    And my repo is left with my original commits
    And Git Town now has the original branch hierarchy
