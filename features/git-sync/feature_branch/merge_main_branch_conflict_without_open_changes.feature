Feature: git sync: resolving merge conflicts between feature and main branch when syncing a feature branch without open changes

  As a developer syncing after having finished a piece of work
  I want to be given an opportunity to resolve differences between my work and the lastest finished features
  So that my work stays stays in sync with the progress of the team, can be easily merged later, and I remain productive.


  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit   | conflicting_file | main content    |
      | feature | local    | conflicting local commit  | conflicting_file | feature content |
    And I run `git sync` while allowing errors


  Scenario: result
    Then I am still on the "feature" branch
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And there is no merge in progress
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                  | FILES            |
      | main    | local and remote | conflicting main commit  | conflicting_file |
      | feature | local            | conflicting local commit | conflicting_file |
    And I still have the following committed files
      | BRANCH  | FILES            | CONTENT         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |


  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature" branch
    And my repo still has a merge in progress


  Scenario Outline: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `<command>`
    Then I am still on the "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILES            |
      | main    | local and remote | conflicting main commit          | conflicting_file |
      | feature | local and remote | Merge branch 'main' into feature |                  |
      |         |                  | conflicting main commit          | conflicting_file |
      |         |                  | conflicting local commit         | conflicting_file |
    And I still have the following committed files
      | BRANCH  | FILES            | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

    Examples:
      | command                                   |
      | git sync --continue                       |
      | git commit --no-edit; git sync --continue |
