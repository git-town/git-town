Feature: git sync: syncing a non-feature branch (with open changes)

  As a developer syncing a non-feature branch
  I want to be able update my ongoing work to include the latest finished features from the rest of the team
  So that our collaboration remains effective.


  Scenario: no conflict
    Given non-feature branch configuration "qa, production"
    And I am on the "qa" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE       | FILE NAME   |
      | qa     | local            | local commit  | local_file  |
      |        | remote           | remote commit | remote_file |
      | main   | local and remote | main commit   | main_file   |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And all branches are now synchronized
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE       | FILES       |
      | qa     | local and remote | local commit  | local_file  |
      |        |                  | remote commit | remote_file |
      | main   | local and remote | main commit   | main_file   |
    And now I have the following committed files
      | BRANCH | FILES                   |
      | qa     | local_file, remote_file |
      | main   | main_file               |
