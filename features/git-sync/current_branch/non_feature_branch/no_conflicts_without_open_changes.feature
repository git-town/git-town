Feature: git sync: syncing a non-feature branch (without open changes)

  (see ./no_conflicts_with_open_changes.feature)


  Background:
    Given non-feature branch configuration "qa, production"
    And I am on the "qa" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE       | FILE NAME   |
      | qa     | local            | local commit  | local_file  |
      |        | remote           | remote commit | remote_file |
      | main   | local and remote | main commit   | main_file   |
    When I run `git sync`


  Scenario: no conflict
    Then it runs the Git commands
      | BRANCH | COMMAND              |
      | qa     | git fetch --prune    |
      | qa     | git rebase origin/qa |
      | qa     | git push             |
      | qa     | git push --tags      |
    And I am still on the "qa" branch
    And all branches are now synchronized
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE       | FILE NAME   |
      | qa     | local and remote | remote commit | remote_file |
      |        |                  | local commit  | local_file  |
      | main   | local and remote | main commit   | main_file   |
