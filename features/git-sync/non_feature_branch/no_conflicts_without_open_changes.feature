Feature: git sync: syncing a non-feature branch (without open changes)

  As a developer working directly off a non-feature branch
  I want to be able update my local work to include the latest finished features from the rest of the team
  So that our work stays in sync, and our collaboration remains productive.


  Scenario: no conflict
    Given non-feature branch configuration "qa, production"
    And I am on the "qa" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE       | FILE NAME   |
      | qa     | local            | local commit  | local_file  |
      |        | remote           | remote commit | remote_file |
      | main   | local and remote | main commit   | main_file   |
    When I run `git sync`
    Then I am still on the "qa" branch
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
