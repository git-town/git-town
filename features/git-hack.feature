Feature: Git Hack

  Scenario: on the main branch
    Given I am on the main branch
    When I run `git hack hot_stuff`
    Then I end up on the "hot_stuff" branch
    And the branch "hot_stuff" has not been pushed to the repository


  Scenario: on the main branch with uncommitted changes
    Given I am on the main branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack hot_stuff`
    Then I end up on the "hot_stuff" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the branch "hot_stuff" has not been pushed to the repository


  Scenario: on the main branch with uncommitted changes and the branch name is taken
    Given I am on the main branch
    And I have a feature branch named "hot_stuff"
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack hot_stuff` while allowing errors
    Then I get the error "A branch named 'hot_stuff' already exists"
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: on a feature branch with uncommitted changes
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message        | file name    |
      | main    | local    | main commit    | main_file    |
      | feature | local    | feature commit | feature_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack hot_stuff`
    Then I end up on the "hot_stuff" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I still have the following commits
      | branch    | location | message        | files        |
      | main      | local    | main commit    | main_file    |
      | feature   | local    | feature commit | feature_file |
      | hot_stuff | local    | main commit    | main_file    |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: with non-pulled updates for the main branch
    Given I am on a feature branch
    And the following commit exists in my repository
      | branch | location | message           | file name        |
      | main   | remote   | new_remote_commit | new_remote_file  |
    When I run `git hack hot_stuff`
    Then I end up on the "hot_stuff" branch
    And I have the following commits
      | branch    | location         | message           | files           |
      | main      | local and remote | new_remote_commit | new_remote_file |
      | hot_stuff | local            | new_remote_commit | new_remote_file |
    And now I have the following committed files
      | branch    | files           |
      | main      | new_remote_file |
      | hot_stuff | new_remote_file |

