Feature: Git Sync: handling conflicting remote branch updates when syncing the main branch without open changes


  Background:
    Given I am on the main branch
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main   | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      | main   | local    | conflicting local commit  | conflicting_file | local conflicting content  |
    And I run `git sync` while allowing errors


  Scenario: result
    Then my repo has a rebase in progress
    And there are abort and continue scripts for "git sync"


  Scenario: aborting
    When I run `git sync --abort`
    Then I am still on the "main" branch
    And there is no rebase in progress
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | BRANCH | LOCATION | MESSAGE                   | FILES              |
      | main   | remote   | conflicting remote commit | conflicting_file   |
      | main   | local    | conflicting local commit  | conflicting_file   |
    And I still have the following committed files
      | BRANCH | FILES              | CONTENT                   |
      | main   | conflicting_file   | local conflicting content |


  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And my repo still has a rebase in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync --continue`
    Then I am still on the "main" branch
    And there are no abort and continue scripts for "git sync" anymore
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILES            |
      | main   | local and remote | conflicting remote commit | conflicting_file |
      | main   | local and remote | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | FILES            | CONTENT          |
      | main   | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run `git rebase --continue`
    When I run `git sync --continue`
    Then I am still on the "main" branch
    And there are no abort and continue scripts for "git sync" anymore
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILES            |
      | main   | local and remote | conflicting remote commit | conflicting_file |
      | main   | local and remote | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | FILES            | CONTENT          |
      | main   | conflicting_file | resolved content |
