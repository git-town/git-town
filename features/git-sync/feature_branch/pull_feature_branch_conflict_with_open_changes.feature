Feature: Git Sync: handling conflicting remote feature branch updates when syncing a feature branch with open changes


  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME          | FILE CONTENT               |
      | feature | remote   | remote conflicting commit | conflicting_file   | remote conflicting content |
      | feature | local    | local conflicting commit  | conflicting_file   | local conflicting content  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git sync` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then I am still on the "feature" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a merge in progress
    And there are abort and continue scripts for "git sync"


  Scenario: aborting
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no merge in progress
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | BRANCH  | LOCATION | MESSAGE                   | FILES              |
      | feature | local    | local conflicting commit  | conflicting_file   |
      | feature | remote   | remote conflicting commit | conflicting_file   |
    And I still have the following committed files
      | BRANCH  | FILES              | CONTENT                   |
      | feature | conflicting_file   | local conflicting content |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo still has a merge in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync --continue`
    Then I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there are no abort and continue scripts for "git sync" anymore
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILES            |
      | feature | local and remote | Merge remote-tracking branch 'origin/feature' into feature |                  |
      | feature | local and remote | remote conflicting commit                                  | conflicting_file |
      | feature | local and remote | local conflicting commit                                   | conflicting_file |
    And now I have the following committed files
      | BRANCH  | FILES              | CONTENT            |
      | feature | conflicting_file   | resolved content   |


  Scenario: continuing after resolving conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    And I run `git commit --no-edit`
    When I run `git sync --continue`
    Then I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there are no abort and continue scripts for "git sync" anymore
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILES            |
      | feature | local and remote | Merge remote-tracking branch 'origin/feature' into feature |                  |
      | feature | local and remote | remote conflicting commit                                  | conflicting_file |
      | feature | local and remote | local conflicting commit                                   | conflicting_file |
    And now I have the following committed files
      | BRANCH  | FILES              | CONTENT            |
      | feature | conflicting_file   | resolved content   |
