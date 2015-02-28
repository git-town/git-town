Feature: git extract: resolving conflicts between main branch and extracted commits (with open changes)

  As a developer extracting a commit that conflicts with the main branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE         | FILE NAME        | FILE CONTENT     |
      | main    | local    | main commit     | conflicting_file | main content     |
      | feature | local    | feature commit  | feature_file     |                  |
      |         |          | refactor commit | conflicting_file | refactor content |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract refactor` with the last commit sha


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                               |
      | feature  | git fetch --prune                     |
      | feature  | git stash -u                          |
      | feature  | git checkout main                     |
      | main     | git rebase origin/main                |
      | main     | git push                              |
      | main     | git checkout -b refactor main         |
      | refactor | git cherry-pick [SHA:refactor commit] |
    And I get the error
      """
      To abort, run "git extract --abort".
      To continue after you have resolved the conflicts, run "git extract --continue".
      """
    And I end up on the "refactor" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a cherry-pick in progress


  Scenario: aborting
    When I run `git extract --abort`
    Then it runs the Git commands
      | BRANCH   | COMMAND                 |
      | refactor | git cherry-pick --abort |
      | refactor | git checkout main       |
      | main     | git branch -d refactor  |
      | main     | git checkout feature    |
      | feature  | git stash pop           |
    And I end up on the "feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no "refactor" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE         | FILE NAME        |
      | main    | local and remote | main commit     | conflicting_file |
      | feature | local            | feature commit  | feature_file     |
      |         |                  | refactor commit | conflicting_file |
    And my repo has no cherry-pick in progress


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git extract --continue`
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git extract"
    And I am still on the "refactor" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a cherry-pick in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git extract --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND                     |
      | refactor | git commit --no-edit        |
      | refactor | git push -u origin refactor |
      | refactor | git stash pop               |
    And I end up on the "refactor" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH   | LOCATION         | MESSAGE         | FILE NAME        |
      | main     | local and remote | main commit     | conflicting_file |
      | feature  | local            | feature commit  | feature_file     |
      |          |                  | refactor commit | conflicting_file |
      | refactor | local and remote | main commit     | conflicting_file |
      |          |                  | refactor commit | conflicting_file |


  Scenario: continuing after resolving the conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git extract --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND                     |
      | refactor | git push -u origin refactor |
      | refactor | git stash pop               |
    And I end up on the "refactor" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH   | LOCATION         | MESSAGE         | FILE NAME        |
      | main     | local and remote | main commit     | conflicting_file |
      | feature  | local            | feature commit  | feature_file     |
      |          |                  | refactor commit | conflicting_file |
      | refactor | local and remote | main commit     | conflicting_file |
      |          |                  | refactor commit | conflicting_file |
