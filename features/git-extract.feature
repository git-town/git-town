Feature: Git Extract

  Scenario: on the main branch
    Given I am on the main branch
    When I run `git extract` while allowing errors
    Then I am still on the "main" branch


  Scenario: on a feature branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message            | file name        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      | feature | local    | refactor commit    | refactor_file    |
    When I run `git extract refactor` with the last commit sha as an argument
    Then I end up on the "refactor" branch
    And all branches are now synchronized
    And I have the following commits
      | branch   | location         | message            | files            |
      | main     | local and remote | remote main commit | remote_main_file |
      | feature  | local            | feature commit     | feature_file     |
      | feature  | local            | refactor commit    | refactor_file    |
      | refactor | local and remote | remote main commit | remote_main_file |
      | refactor | local and remote | refactor commit    | refactor_file    |
    And now I have the following committed files
      | branch   | files                           |
      | main     | remote_main_file                |
      | feature  | feature_file, refactor_file     |
      | refactor | remote_main_file, refactor_file |


  Scenario: user aborts after merge conflict during cherry-picking
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message            | file name        | file content    |
      | main    | local    | conflicting commit | conflicting_file | main content    |
      | feature | local    | conflicting commit | conflicting_file | feature content |
    When I run `git extract refactor` with the last commit sha as an argument while allowing errors
    Then I end up on the "refactor" branch
    And my repo has a cherry-pick in progress
    And there is an abort script for "git extract"
    When I run `git extract --abort`
    Then I end up on the "feature" branch
    And I have the following commits
      | branch   | location | message            | files            |
      | main     | local    | conflicting commit | conflicting_file |
      | feature  | local    | conflicting commit | conflicting_file |
    And my repo has no cherry-pick in progress
    And there is no abort script for "git extract" anymore


  Scenario: user aborts after merge conflict during main branch pulling
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message                   | file name        | file content   |
      | main    | remote   | conflicting remote commit | conflicting_file | remote content |
      | main    | local    | conflicting local commit  | conflicting_file | local content  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract refactor` with the last commit sha as an argument while allowing errors
    Then I get the error "ERROR WHILE PULLING THE MAIN BRANCH"
    And my repo has a rebase in progress
    And there is an abort script for "git extract"
    And I don't have an uncommitted file with name: "uncommitted"
    When I run `git extract --abort`
    Then I end up on my feature branch
    And I have the feature branches "feature"
    And I have the following commits
      | branch | location | message                   | files            |
      | main   | remote   | conflicting remote commit | conflicting_file |
      | main   | local    | conflicting local commit  | conflicting_file |
    And there is no rebase in progress
    And there is no abort script for "git extract" anymore
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"

