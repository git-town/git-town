Feature: sync the main branch

  Background:
    Given I am on the "main" branch
    And the commits
      | LOCATION | MESSAGE       | FILE NAME   |
      | local    | local commit  | local_file  |
      | remote   | remote commit | remote_file |
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/main   |
      |        | git push                 |
      |        | git push --tags          |
      |        | git stash pop            |
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, remote | remote commit |
      |        |               | local commit  |
