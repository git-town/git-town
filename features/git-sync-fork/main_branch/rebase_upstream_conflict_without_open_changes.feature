Feature: git-sync-fork: handling rebase conflicts between main branch and its remote with open changes

  Background:
    Given my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         | FILE NAME        | FILE CONTENT     |
      | main   | upstream | upstream commit | conflicting_file | upstream content |
      |        | local    | local commit    | conflicting_file | local content    |
    And I am on the "main" branch
    When I run `git sync-fork` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                  |
      | main   | git fetch upstream       |
      | main   | git rebase upstream/main |
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git sync-fork --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND            |
      | HEAD   | git rebase --abort |
    And I end up on the "main" branch
    And there is no rebase in progress
    And I still have the following commits
      | BRANCH | LOCATION | MESSAGE         | FILES            |
      | main   | upstream | upstream commit | conflicting_file |
      |        | local    | local commit    | conflicting_file |
    And I still have the following committed files
      | BRANCH | FILES            | CONTENT       |
      | main   | conflicting_file | local content |


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync-fork --continue`
    Then it runs the Git commands
      | BRANCH | COMMAND               |
      | HEAD   | git rebase --continue |
      | main   | git push              |
    And I end up on the "main" branch
    And I still have the following commits
      | BRANCH | LOCATION                    | MESSAGE         | FILES            |
      | main   | local, remote, and upstream | upstream commit | conflicting_file |
      | main   | local, remote               | local commit    | conflicting_file |
    And now I have the following committed files
      | BRANCH | FILES            | CONTENT          |
      | main   | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git sync-fork --continue`
    Then it runs the Git commands
      | BRANCH | COMMAND  |
      | main   | git push |
    And I end up on the "main" branch
    And I still have the following commits
      | BRANCH | LOCATION                    | MESSAGE         | FILES            |
      | main   | local, remote, and upstream | upstream commit | conflicting_file |
      | main   | local, remote               | local commit    | conflicting_file |
    And now I have the following committed files
      | BRANCH | FILES            | CONTENT          |
      | main   | conflicting_file | resolved content |
