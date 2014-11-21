Feature: Git Sync: handling conflicting remote main branch updates when syncing a feature branch with open changes


  Background:
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message                   | file name        | file content               |
      | main    | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      | main    | local    | conflicting local commit  | conflicting_file | local conflicting content  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git sync` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then my repo has a rebase in progress
    And there are abort and continue scripts for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | branch  | location | message                   | files              |
      | main    | remote   | conflicting remote commit | conflicting_file   |
      | main    | local    | conflicting local commit  | conflicting_file   |
    And I still have the following committed files
      | branch | files              | content                   |
      | main   | conflicting_file   | local conflicting content |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And my repo still has a rebase in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
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


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run `git rebase --continue`
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
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
