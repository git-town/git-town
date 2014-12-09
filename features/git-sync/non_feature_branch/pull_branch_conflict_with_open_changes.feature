Feature: Git Sync: handling conflicting remote branch updates when syncing a non-feature branch with open changes

  Background:
    Given non-feature branch configuration "qa, production"
    And I am on the "qa" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | qa     | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      |        | local    | conflicting local commit  | conflicting_file | local conflicting content  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git sync` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then my repo has a rebase in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git sync --abort`
    Then I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And there are no abort and continue scripts for "git sync"
    And I still have the following commits
      | BRANCH | LOCATION | MESSAGE                   | FILES              |
      | qa     | remote   | conflicting remote commit | conflicting_file   |
      |        | local    | conflicting local commit  | conflicting_file   |
    And I still have the following committed files
      | BRANCH | FILES              | CONTENT                   |
      | qa     | conflicting_file   | local conflicting content |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo still has a rebase in progress


  Scenario Outline: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `<command>`
    Then I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there are no abort and continue scripts for "git sync"
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILES            |
      | qa     | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | FILES            | CONTENT          |
      | qa     | conflicting_file | resolved content |

    Examples:
      | command                                    |
      | git sync --continue                        |
      | git rebase --continue; git sync --continue |
