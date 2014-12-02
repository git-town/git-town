Feature: git-hack handling conflicting remote main branch updates with open changes

  Background:
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | remote   | conflicting remote commit | conflicting_file | remote content |
      |        | local    | conflicting local commit  | conflicting_file | local content  |
    And I am on the "feature" branch
    When I run `git hack other_feature` while allowing errors


  Scenario: result
    Then my repo has a rebase in progress


  Scenario: aborting
    When I run `git hack --abort`
    Then I end up on the "feature" branch
    And there is no rebase in progress
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE                   | FILES            |
      | main   | remote   | conflicting remote commit | conflicting_file |
      |        | local    | conflicting local commit  | conflicting_file |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git hack --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git hack"
    And my repo still has a rebase in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git hack --continue`
    Then I end up on the "other_feature" branch
    And now I have the following commits
      | BRANCH        | LOCATION         | MESSAGE                   | FILES            |
      | main          | local and remote | conflicting remote commit | conflicting_file |
      | main          | local and remote | conflicting local commit  | conflicting_file |
      | other_feature | local            | conflicting remote commit | conflicting_file |
      | other_feature | local            | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH        | FILES            | CONTENT          |
      | main          | conflicting_file | resolved content |
      | other_feature | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run `git rebase --continue`
    When I run `git hack --continue`
    Then I end up on the "other_feature" branch
    And now I have the following commits
      | BRANCH        | LOCATION         | MESSAGE                   | FILES            |
      | main          | local and remote | conflicting remote commit | conflicting_file |
      | main          | local and remote | conflicting local commit  | conflicting_file |
      | other_feature | local            | conflicting remote commit | conflicting_file |
      | other_feature | local            | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH        | FILES            | CONTENT          |
      | main          | conflicting_file | resolved content |
      | other_feature | conflicting_file | resolved content |
