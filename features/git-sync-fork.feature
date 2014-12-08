Feature: Git Sync-Fork

  Scenario: on the main branch with an upstream commit
    Given I am on the main branch
    And my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         | FILE NAME     |
      | main   | upstream | upstream commit | upstream_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-fork`
    Then I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH | LOCATION                    | MESSAGE         | FILES         |
      | main   | local, remote, and upstream | upstream commit | upstream_file |
    And now I have the following committed files
      | BRANCH | FILES         |
      | main   | upstream_file |


  Scenario: on a feature branch with upstream commit in main branch
    Given I am on a feature branch
    And my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         | FILE NAME     |
      | main   | upstream | upstream commit | upstream_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-fork`
    Then I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH | LOCATION                | MESSAGE         | FILES         |
      | main   | local, remote, upstream | upstream commit | upstream_file |
    And now I have the following committed files
      | BRANCH | FILES         |
      | main   | upstream_file |


  Scenario: user aborts after conflict while pulling upstream updates
    Given I am on a feature branch
    And my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         | FILE NAME        | FILE CONTENT     |
      | main   | upstream | upstream commit | conflicting_file | upstream content |
      |        | local    | local commit    | conflicting_file | local content    |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-fork` while allowing errors
    Then my repo has a rebase in progress
    And there is an abort script for "git sync-fork"
    And I don't have an uncommitted file with name: "uncommitted"

    When I run `git sync-fork --abort`
    Then I end up on my feature branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And there is no abort script for "git sync-fork" anymore
    And I still have the following commits
      | BRANCH | LOCATION | MESSAGE         | FILES            |
      | main   | upstream | upstream commit | conflicting_file |
      |        | local    | local commit    | conflicting_file |
    And I still have the following committed files
      | BRANCH | FILES            | CONTENT       |
      | main   | conflicting_file | local content |


  Scenario: without upstream
    When I run `git sync-fork` while allowing errors
    Then I get the error "Please add a remote 'upstream'"
