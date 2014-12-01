Feature: git-hack handling conflicting remote main branch updates with open changes

  Background:
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | branch | location | message                   | file name        | file content   |
      | main   | remote   | conflicting remote commit | conflicting_file | remote content |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack other_feature` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then my repo has a rebase in progress
    And there are abort and continue script for "git hack"
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git hack --abort`
    Then I end up on the "feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And there is no abort script for "git hack" anymore
    And I have the following commits
      | branch | location | message                   | files            |
      | main   | remote   | conflicting remote commit | conflicting_file |
      | main   | local    | conflicting local commit  | conflicting_file |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git hack --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git hack"
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo still has a rebase in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git hack --continue`
    Then I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there are no abort and continue scripts for "git hack" anymore
    And now I have the following commits
      | branch        | location         | message                   | files            |
      | main          | local and remote | conflicting remote commit | conflicting_file |
      | main          | local and remote | conflicting local commit  | conflicting_file |
      | other_feature | local            | conflicting remote commit | conflicting_file |
      | other_feature | local            | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | branch        | files            | content          |
      | main          | conflicting_file | resolved content |
      | other_feature | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run `git rebase --continue`
    When I run `git hack --continue`
    Then I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there are no abort and continue scripts for "git hack" anymore
    And now I have the following commits
      | branch        | location         | message                   | files            |
      | main          | local and remote | conflicting remote commit | conflicting_file |
      | main          | local and remote | conflicting local commit  | conflicting_file |
      | other_feature | local            | conflicting remote commit | conflicting_file |
      | other_feature | local            | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | branch        | files            | content          |
      | main          | conflicting_file | resolved content |
      | other_feature | conflicting_file | resolved content |
