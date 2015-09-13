Feature: git extract: resolving conflicts between main branch and extracted commits (without remote origin)

  (see ../with_remote_origin/with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE         | FILE NAME        | FILE CONTENT     |
      | main    | local    | main commit     | conflicting_file | main content     |
      | feature | local    | feature commit  | feature_file     |                  |
      |         |          | refactor commit | conflicting_file | refactor content |
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git extract refactor` with the last commit sha


  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                                      |
      | feature  | git stash -u                                 |
      |          | git checkout -b refactor main                |
      | refactor | git cherry-pick <%= sha 'refactor commit' %> |
    And I get the error
      """
      To abort, run "git extract --abort".
      To continue after you have resolved the conflicts, run "git extract --continue".
      """
    And I end up on the "refactor" branch
    And my uncommitted file is stashed
    And my repo has a cherry-pick in progress


  Scenario: aborting
    When I run `git extract --abort`
    Then it runs the commands
      | BRANCH   | COMMAND                 |
      | refactor | git cherry-pick --abort |
      |          | git checkout feature    |
      | feature  | git branch -D refactor  |
      |          | git stash pop           |
    And I end up on the "feature" branch
    And I again have my uncommitted file
    And there is no "refactor" branch
    And I am left with my original commits
    And my repo has no cherry-pick in progress


  Scenario: continuing without resolving the conflicts
    When I run `git extract --continue`
    Then it runs no commands
    And I get the error "You must resolve the conflicts before continuing the git extract"
    And I am still on the "refactor" branch
    And my uncommitted file is stashed
    And my repo has a cherry-pick in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git extract --continue`
    Then it runs the commands
      | BRANCH   | COMMAND              |
      | refactor | git commit --no-edit |
      |          | git stash pop        |
    And I end up on the "refactor" branch
    And I again have my uncommitted file
    And now I have the following commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | local    | main commit     |
      | feature  | local    | feature commit  |
      |          |          | refactor commit |
      | refactor | local    | main commit     |
      |          |          | refactor commit |
    And now I have the following committed files
      | BRANCH   | NAME             | CONTENT          |
      | main     | conflicting_file | main content     |
      | feature  | conflicting_file | refactor content |
      | feature  | feature_file     |                  |
      | refactor | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git extract --continue`
    Then it runs the commands
      | BRANCH   | COMMAND       |
      | refactor | git stash pop |
    And I end up on the "refactor" branch
    And I again have my uncommitted file
    And now I have the following commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | local    | main commit     |
      | feature  | local    | feature commit  |
      |          |          | refactor commit |
      | refactor | local    | main commit     |
      |          |          | refactor commit |
    And now I have the following committed files
      | BRANCH   | NAME             | CONTENT          |
      | main     | conflicting_file | main content     |
      | feature  | conflicting_file | refactor content |
      | feature  | feature_file     |                  |
      | refactor | conflicting_file | resolved content |
