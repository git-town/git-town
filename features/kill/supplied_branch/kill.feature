Feature: delete another than the current branch

  Background:
    Given my repo has the feature branches "good" and "dead"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE                              | FILE NAME        |
      | main   | local, remote | conflicting with uncommitted changes | conflicting_file |
      | dead   | local, remote | dead-end commit                      | file             |
      | good   | local, remote | good commit                          | file             |
    And I am on the "good" branch
    And my workspace has an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town kill dead"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | good   | git fetch --prune --tags |
      |        | git push origin :dead    |
      |        | git branch -D dead       |
    And I am still on the "good" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY    | BRANCHES   |
      | local, remote | main, good |
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE                              |
      | main   | local, remote | conflicting with uncommitted changes |
      | good   | local, remote | good commit                          |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | good   | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | good   | git branch dead {{ sha 'dead-end commit' }} |
      |        | git push -u origin dead                     |
    And I am still on the "good" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
    And my repo now has its initial branches and branch hierarchy
