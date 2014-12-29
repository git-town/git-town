Feature: git extract: resolving conflicts with main branch (without open changes)

  (see ./cherry_pick_conflict_with_open_changes.feature)


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
    Then it runs the Git commands
      | BRANCH   | COMMAND                               |
      | feature  | git fetch --prune                     |
      | feature  | git checkout main                     |
      | main     | git rebase origin/main                |
      | main     | git push                              |
      | main     | git checkout -b refactor main         |
      | refactor | git cherry-pick [SHA:refactor commit] |
    And I end up on the "refactor" branch
    And my repo has a cherry-pick in progress


  Scenario: aborting
    When I run `git extract --abort`
    Then it runs the Git commands
      | BRANCH   | COMMAND                 |
      | refactor | git cherry-pick --abort |
      | refactor | git checkout main       |
      | main     | git branch -d refactor  |
      | main     | git checkout feature    |
    And I end up on the "feature" branch
    And there is no "refactor" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE         | FILES            |
      | main    | local and remote | main commit     | conflicting_file |
      | feature | local            | feature commit  | feature_file     |
      |         |                  | refactor commit | conflicting_file |
    And my repo has no cherry-pick in progress


  Scenario: continuing without resolving conflicts
    When I run `git extract --continue` while allowing errors
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git extract"
    And I am still on the "refactor" branch
    And my repo has a cherry-pick in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git extract --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND                     |
      | refactor | git commit --no-edit        |
      | refactor | git push -u origin refactor |
    And I end up on the "refactor" branch
    And now I have the following commits
      | BRANCH   | LOCATION         | MESSAGE         | FILES            |
      | main     | local and remote | main commit     | conflicting_file |
      | feature  | local            | feature commit  | feature_file     |
      |          |                  | refactor commit | conflicting_file |
      | refactor | local and remote | main commit     | conflicting_file |
      |          |                  | refactor commit | conflicting_file |


  Scenario: continuing after resolving conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git extract --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND                     |
      | refactor | git push -u origin refactor |
    And I end up on the "refactor" branch
    And now I have the following commits
      | BRANCH   | LOCATION         | MESSAGE         | FILES            |
      | main     | local and remote | main commit     | conflicting_file |
      | feature  | local            | feature commit  | feature_file     |
      |          |                  | refactor commit | conflicting_file |
      | refactor | local and remote | main commit     | conflicting_file |
      |          |                  | refactor commit | conflicting_file |
