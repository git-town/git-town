Feature: Git Sync: handling conflicting remote feature branch updates when syncing a feature branch without open changes


  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
      |         | local    | local conflicting commit  | conflicting_file | local conflicting content  |
    And I run `git sync` while allowing errors


  Scenario: result
    Then I am still on the "feature" branch
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And there is no merge in progress
    And I still have the following commits
      | BRANCH  | LOCATION | MESSAGE                   | FILES            |
      | feature | local    | local conflicting commit  | conflicting_file |
      |         | remote   | remote conflicting commit | conflicting_file |
    And I still have the following committed files
      | BRANCH  | FILES            | CONTENT                   |
      | feature | conflicting_file | local conflicting content |


  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature" branch
    And my repo still has a merge in progress


  Scenario Outline: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `<COMMAND>`
    Then I am still on the "feature" branch
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILES            |
      | feature | local and remote | Merge remote-tracking branch 'origin/feature' into feature |                  |
      |         |                  | remote conflicting commit                                  | conflicting_file |
      |         |                  | local conflicting commit                                   | conflicting_file |
    And now I have the following committed files
      | BRANCH  | FILES            | CONTENT          |
      | feature | conflicting_file | resolved content |

    Examples:
      | COMMAND                                   |
      | git sync --continue                       |
      | git commit --no-edit; git sync --continue |
