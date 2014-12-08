Feature: Git Sync: handling merge conflicts between feature and main branch when syncing a feature branch with open changes


  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit   | conflicting_file | main content    |
      | feature | local    | conflicting local commit  | conflicting_file | feature content |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git sync` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then I am still on the "feature" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no merge in progress
    And there are no abort and continue scripts for "git sync"
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                  | FILES            |
      | main    | local and remote | conflicting main commit  | conflicting_file |
      | feature | local            | conflicting local commit | conflicting_file |
    And I still have the following committed files
      | BRANCH  | FILES            | CONTENT         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo still has a merge in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there are no abort and continue scripts for "git sync"
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILES            |
      | main    | local and remote | conflicting main commit          | conflicting_file |
      | feature | local and remote | Merge branch 'main' into feature |                  |
      | feature | local and remote | conflicting main commit          | conflicting_file |
      | feature | local and remote | conflicting local commit         | conflicting_file |
    And I still have the following committed files
      | BRANCH  | FILES            | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    And I run `git commit --no-edit`
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there are no abort and continue scripts for "git sync"
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILES            |
      | main    | local and remote | conflicting main commit          | conflicting_file |
      | feature | local and remote | Merge branch 'main' into feature |                  |
      | feature | local and remote | conflicting main commit          | conflicting_file |
      | feature | local and remote | conflicting local commit         | conflicting_file |
    And I still have the following committed files
      | BRANCH  | FILES            | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
