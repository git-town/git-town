Feature: git-town sync: syncing the current perennial branch

  As a developer syncing a perennial branch
  I want to be able update my ongoing work to include the latest finished features from the rest of the team
  So that our collaboration remains effective.


  Background:
    Given my repository has the perennial branches "production" and "qa"
    And the following commits exist in my repository
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME   |
      | qa     | local         | local commit  | local_file  |
      |        | remote        | remote commit | remote_file |
      | main   | local, remote | main commit   | main_file   |
    And I am on the "qa" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"


  Scenario: no conflict
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | qa     | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/qa     |
      |        | git push                 |
      |        | git push --tags          |
      |        | git stash pop            |
    And I am still on the "qa" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
    And my repository now has the following commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME   |
      | main   | local, remote | main commit   | main_file   |
      | qa     | local, remote | remote commit | remote_file |
      |        |               | local commit  | local_file  |
