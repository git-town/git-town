Feature: git town-sync: syncing the main branch

  As a developer syncing the main branch
  I want to be able update my ongoing work to include the latest finished features from the rest of the team
  So that our collaboration remains effective.


  Background:
    Given I am on the "main" branch
    And the following commits exist in my repository
      | LOCATION | MESSAGE       | FILE NAME   |
      | local    | local commit  | local_file  |
      | remote   | remote commit | remote_file |
    And I have an uncommitted file
    When I run `git town-sync`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                |
      | main   | git fetch --prune      |
      |        | git stash -a           |
      |        | git rebase origin/main |
      |        | git push               |
      |        | git push --tags        |
      |        | git stash pop          |
    And I am still on the "main" branch
    And I still have my uncommitted file
    And all branches are now synchronized
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE       | FILE NAME   |
      | main   | local and remote | remote commit | remote_file |
      |        |                  | local commit  | local_file  |
