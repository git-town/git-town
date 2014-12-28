Feature: git extract: allows to resolve conflicting remote main branch updates (with open changes)

  As a developer extracting a commit when the main branch has conflicting local and remote updates
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given I am on a feature branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main    | remote   | conflicting remote commit | conflicting_file | remote content |
      |         | local    | conflicting local commit  | conflicting_file | local content  |
      | feature | local    | feature commit            | feature_file     |                |
      |         |          | refactor commit           | refactor_file    |                |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract refactor` with the last commit sha while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then my repo has a rebase in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git extract --abort`
    Then I end up on my feature branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no "refactor" branch
    And I have the following commits
      | BRANCH  | LOCATION | MESSAGE                   | FILES            |
      | main    | remote   | conflicting remote commit | conflicting_file |
      |         | local    | conflicting local commit  | conflicting_file |
      | feature | local    | feature commit            | feature_file     |
      |         |          | refactor commit           | refactor_file    |
    And there is no rebase in progress


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git extract --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git extract"
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a rebase in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git extract --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND                                  |
      | HEAD     | git rebase --continue                    |
      | main     | git push                                 |
      | main     | git checkout -b refactor main            |
      | refactor | git cherry-pick [["feature" branch SHA]] |
      | refactor | git push -u origin refactor              |
      | refactor | git stash pop                            |
    And I end up on the "refactor" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH   | LOCATION         | MESSAGE                   | FILES            |
      | main     | local and remote | conflicting remote commit | conflicting_file |
      |          |                  | conflicting local commit  | conflicting_file |
      | feature  | local            | feature commit            | feature_file     |
      |          |                  | refactor commit           | refactor_file    |
      | refactor | local and remote | conflicting remote commit | conflicting_file |
      |          |                  | conflicting local commit  | conflicting_file |
      |          |                  | refactor commit           | refactor_file    |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git extract --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND                                  |
      | main     | git push                                 |
      | main     | git checkout -b refactor main            |
      | refactor | git cherry-pick [["feature" branch SHA]] |
      | refactor | git push -u origin refactor              |
      | refactor | git stash pop                            |
    And I end up on the "refactor" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH   | LOCATION         | MESSAGE                   | FILES            |
      | main     | local and remote | conflicting remote commit | conflicting_file |
      |          |                  | conflicting local commit  | conflicting_file |
      | feature  | local            | feature commit            | feature_file     |
      |          |                  | refactor commit           | refactor_file    |
      | refactor | local and remote | conflicting remote commit | conflicting_file |
      |          |                  | conflicting local commit  | conflicting_file |
      |          |                  | refactor commit           | refactor_file    |
