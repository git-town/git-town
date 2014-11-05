Feature: git-sync on a feature branch with open changes (conflict when pulling the feature branch)

  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | branch  | location | message                   | file name          | file content               |
      | feature | remote   | remote conflicting commit | conflicting_file   | remote conflicting content |
      | feature | local    | local conflicting commit  | conflicting_file   | local conflicting content  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git sync` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then I am still on the "feature" branch
    And my repo has a merge in progress
    And there are abort and continue scripts for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And there is no merge in progress
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | branch  | location | message                   | files              |
      | feature | local    | local conflicting commit  | conflicting_file   |
      | feature | remote   | remote conflicting commit | conflicting_file   |
    And I still have the following committed files
      | branch  | files              | content                   |
      | feature | conflicting_file   | local conflicting content |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then I get the error "You must resolve the conflicts and commit your changes before continuing the git sync."
    And I am still on the "feature" branch
    And my repo still has a merge in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: continuing after resolving conflicts
    When I successfully finish the merge by resolving the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And there are no abort and continue scripts for "git sync" anymore
    And now I have the following commits
      | branch  | location         | message                   | files              |
      | feature | local and remote | remote conflicting commit | conflicting_file   |
      | feature | local and remote | local conflicting commit  | conflicting_file   |
    And now I have the following committed files
      | branch  | files              | content            |
      | feature | conflicting_file   | resolved content   |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
