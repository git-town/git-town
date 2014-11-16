Feature: Git Sync-Fork

  Scenario: on the main branch with an upstream commit
    Given I am on the main branch
    And my repo has an upstream repo
    And the following commits exist in my repository
      | branch | location | message         | file name     |
      | main   | upstream | upstream commit | upstream_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-fork`
    Then I am still on the "main" branch
    And I see the following commits
      | branch | location                    | message         | files         |
      | main   | local, remote, and upstream | upstream commit | upstream_file |
    And now I have the following committed files
      | branch | files         |
      | main   | upstream_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: on a feature branch with upstream commit in main branch
    Given I am on a feature branch
    And my repo has an upstream repo
    And the following commits exist in my repository
      | branch | location | message         | file name     |
      | main   | upstream | upstream commit | upstream_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-fork`
    Then I am still on the "feature" branch
    And I see the following commits
      | branch | location         | message         | files         |
      | main   | local and remote | upstream commit | upstream_file |
    And now I have the following committed files
      | branch | files         |
      | main   | upstream_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user aborts after conflict while pulling upstream updates
    Given I am on a feature branch
    And my repo has an upstream repo
    And the following commits exist in my repository
      | branch | location | message         | file name        | file content     |
      | main   | upstream | upstream commit | conflicting_file | upstream content |
      | main   | local    | local commit    | conflicting_file | local content    |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-fork` while allowing errors
    Then my repo has a rebase in progress
    And there is an abort script for "git sync-fork"
    And I don't have an uncommitted file with name: "uncommitted"

    When I run `git sync-fork --abort`
    Then I end up on my feature branch
    And there is no rebase in progress
    And there is no abort script for "git sync-fork" anymore
    And I still have the following commits
      | branch | location | message      | files            |
      | main   | local    | local commit | conflicting_file |
    And I still have the following committed files
      | branch | files            | content       |
      | main   | conflicting_file | local content |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: without upstream
    When I run `git sync-fork` while allowing errors
    Then I get the error "Please add a remote 'upstream'"
