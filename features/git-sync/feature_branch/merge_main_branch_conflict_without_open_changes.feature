Feature: Git Sync: handling merge conflicts between feature and main branch when syncing a feature branch without open changes


  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | branch  | location | message                   | file name        | file content    |
      | main    | local    | conflicting main commit   | conflicting_file | main content    |
      | feature | local    | conflicting local commit  | conflicting_file | feature content |
    And I run `git sync` while allowing errors


  Scenario: result
    Then I am still on the "feature" branch
    And my repo has a merge in progress
    And there are abort and continue scripts for "git sync"


  Scenario: aborting
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And there is no merge in progress
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | branch  | location | message                  | files            |
      | main    | local    | conflicting main commit  | conflicting_file |
      | feature | local    | conflicting local commit | conflicting_file |
    And I still have the following committed files
      | branch  | files            | content         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |


  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature" branch
    And my repo still has a merge in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | branch  | location         | message                          | files            |
      | main    | local            | conflicting main commit          | conflicting_file |
      | feature | local and remote | Merge branch 'main' into feature |                  |
      | feature | local and remote | conflicting main commit          | conflicting_file |
      | feature | local and remote | conflicting local commit         | conflicting_file |
    And I still have the following committed files
      | branch  | files            | content          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    And I run `git commit --no-edit`
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | branch  | location         | message                          | files            |
      | main    | local            | conflicting main commit          | conflicting_file |
      | feature | local and remote | Merge branch 'main' into feature |                  |
      | feature | local and remote | conflicting main commit          | conflicting_file |
      | feature | local and remote | conflicting local commit         | conflicting_file |
    And I still have the following committed files
      | branch  | files            | content          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
