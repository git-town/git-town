Feature: git-extract handling cherry-pick conflicts without open changes

  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE         | FILE NAME        | FILE CONTENT     |
      | main    | local    | main commit     | conflicting_file | main content     |
      | feature | local    | feature commit  | feature_file     |                  |
      |         |          | refactor commit | conflicting_file | refactor content |
    And I am on the "feature" branch
    When I run `git extract refactor` with the last commit sha while allowing errors


  Scenario: result
    Then I end up on the "refactor" branch
    And my repo has a cherry-pick in progress


  Scenario: aborting
    When I run `git extract --abort`
    Then I end up on the "feature" branch
    And there is no "refactor" branch
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE         | FILES            |
      | main     | local and remote | main commit     | conflicting_file |
      | feature  | local            | feature commit  | feature_file     |
      |          |                  | refactor commit | conflicting_file |
    And my repo has no cherry-pick in progress


  Scenario: continuing without resolving conflicts
    When I run `git extract --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git extract"
    And I am still on the "refactor" branch
    And my repo has a cherry-pick in progress


  Scenario Outline: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `<command>`
    Then I end up on the "refactor" branch
    And now I have the following commits
      | BRANCH   | LOCATION         | MESSAGE         | FILES            |
      | main     | local and remote | main commit     | conflicting_file |
      | feature  | local            | feature commit  | feature_file     |
      |          |                  | refactor commit | conflicting_file |
      | refactor | local and remote | main commit     | conflicting_file |
      |          |                  | refactor commit | conflicting_file |

    Examples:
      | command                                      |
      | git extract --continue                       |
      | git commit --no-edit; git extract --continue |
