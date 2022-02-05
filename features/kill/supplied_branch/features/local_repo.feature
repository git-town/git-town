Feature: local repository

  Background:
    Given my repo does not have a remote origin
    And my repo has the local feature branches "good" and "other"
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME        |
      | main   | local    | main commit  | conflicting_file |
      | good   | local    | good commit  | file             |
      | other  | local    | other commit | file             |
    And I am on the "good" branch
    And my workspace has an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town kill other"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND             |
      | good   | git add -A          |
      |        | git stash           |
      |        | git branch -D other |
      |        | git stash pop       |
    And I am still on the "good" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES   |
      | local      | main, good |
    And my repo now has the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main commit |
      | good   | local    | good commit |
    And Git Town now knows this branch hierarchy
      | BRANCH | PARENT |
      | good   | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                   |
      | good   | git add -A                                |
      |        | git stash                                 |
      |        | git branch other {{ sha 'other commit' }} |
      |        | git stash pop                             |
    And I am still on the "good" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my initial commits
    And my repo now has its initial branches and branch hierarchy
