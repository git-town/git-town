Feature: git-sync
  on a feature branch
  conflict when pulling the main branch


  Background:
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message                   | file name        | file content               |
      | main    | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      | main    | local    | conflicting local commit  | conflicting_file | local conflicting content  |
    And I run `git sync` while allowing errors


  Scenario: result
    Then my repo has a rebase in progress
    And there are abort and continue scripts for "git sync"


  Scenario: aborting
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And there is no rebase in progress
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | branch  | location | message                   | files              |
      | main    | remote   | conflicting remote commit | conflicting_file   |
      | main    | local    | conflicting local commit  | conflicting_file   |
    And I still have the following committed files
      | branch | files              | content                   |
      | main   | conflicting_file   | local conflicting content |


  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then I get the error "You must resolve the conflicts and commit your changes before continuing the git sync."
    And my repo still has a rebase in progress


  Scenario: continuing after resolving conflicts
    When I successfully finish the rebase by resolving the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And there are no abort and continue scripts for "git sync" anymore
    And now I have the following commits
      | branch  | location         | message                   | files            |
      | main    | local and remote | conflicting remote commit | conflicting_file |
      | main    | local and remote | conflicting local commit  | conflicting_file |
      | feature | local and remote | conflicting remote commit | conflicting_file |
      | feature | local and remote | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | branch  | files            | content          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |
