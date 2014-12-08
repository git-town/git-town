Feature: Git Sync: syncing the main branch with open changes


  Scenario: no conflicts
    Given I am on the main branch
    And the following commits exist in my repository
      | LOCATION | MESSAGE       | FILE NAME   |
      | local    | local commit  | local_file  |
      | remote   | remote commit | remote_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And all branches are now synchronized
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE       | FILES       |
      | main   | local and remote | local commit  | local_file  |
      |        |                  | remote commit | remote_file |
    And now I have the following committed files
      | BRANCH | FILES       |
      | main   | local_file  |
      | main   | remote_file |
