Feature: git sync: syncing the main branch (without open changes)

  (see ./with_open_changes.feature)

  Background:
    Given I am on the "main" branch
    And the following commits exist in my repository
      | LOCATION | MESSAGE       | FILE NAME   |
      | local    | local commit  | local_file  |
      | remote   | remote commit | remote_file |
    When I run `git sync`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                |
      | main   | git fetch --prune      |
      | main   | git rebase origin/main |
      | main   | git push               |
      | main   | git push --tags        |
    And I am still on the "main" branch
    And all branches are now synchronized
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE       | FILE NAME   |
      | main   | local and remote | remote commit | remote_file |
      |        |                  | local commit  | local_file  |
