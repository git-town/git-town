Feature: git sync --all: handling merge conflicts between feature branch and main branch (with open changes and without remote origin)

  Background:
    Given I have feature branches named "feature1" and "feature2"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH   | LOCATION | MESSAGE         | FILE NAME        | FILE CONTENT     |
      | main     | local    | main commit     | conflicting_file | main content     |
      | feature1 | local    | feature1 commit | conflicting_file | feature1 content |
      | feature2 | local    | feature2 commit | feature2_file    | feature2 content |
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync --all`


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                  |
      | main     | git stash -u             |
      | main     | git checkout feature1    |
      | feature1 | git merge --no-edit main |
    And I get the error
      """
      To abort, run "git sync --abort".
      To continue after you have resolved the conflicts, run "git sync --continue".
      To skip the sync of the 'feature1' branch, run "git sync --skip".
      """
    And I end up on the "feature1" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH   | COMMAND           |
      | feature1 | git merge --abort |
      | feature1 | git checkout main |
      | main     | git stash pop     |
    And I end up on the "main" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH   | LOCATION | MESSAGE         | FILE NAME        |
      | main     | local    | main commit     | conflicting_file |
      | feature1 | local    | feature1 commit | conflicting_file |
      | feature2 | local    | feature2 commit | feature2_file    |


  Scenario: skipping
    When I run `git sync --skip`
    Then it runs the Git commands
      | BRANCH   | COMMAND                  |
      | feature1 | git merge --abort        |
      | feature1 | git checkout feature2    |
      | feature2 | git merge --no-edit main |
      | feature2 | git checkout main        |
      | main     | git stash pop            |
    And I end up on the "main" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH   | LOCATION | MESSAGE                           | FILE NAME        |
      | main     | local    | main commit                       | conflicting_file |
      | feature1 | local    | feature1 commit                   | conflicting_file |
      | feature2 | local    | feature2 commit                   | feature2_file    |
      |          |          | main commit                       | conflicting_file |
      |          |          | Merge branch 'main' into feature2 |                  |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving the conflicts
    When I run `git sync --continue`
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature1" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND                  |
      | feature1 | git commit --no-edit     |
      | feature1 | git checkout feature2    |
      | feature2 | git merge --no-edit main |
      | feature2 | git checkout main        |
      | main     | git stash pop            |
    And I end up on the "main" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH   | LOCATION | MESSAGE                           | FILE NAME        |
      | main     | local    | main commit                       | conflicting_file |
      | feature1 | local    | feature1 commit                   | conflicting_file |
      |          |          | main commit                       | conflicting_file |
      |          |          | Merge branch 'main' into feature1 |                  |
      | feature2 | local    | feature2 commit                   | feature2_file    |
      |          |          | main commit                       | conflicting_file |
      |          |          | Merge branch 'main' into feature2 |                  |


  Scenario: continuing after resolving the conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    And I run `git commit --no-edit; git sync --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND                  |
      | feature1 | git checkout feature2    |
      | feature2 | git merge --no-edit main |
      | feature2 | git checkout main        |
      | main     | git stash pop            |
    And I end up on the "main" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH   | LOCATION | MESSAGE                           | FILE NAME        |
      | main     | local    | main commit                       | conflicting_file |
      | feature1 | local    | feature1 commit                   | conflicting_file |
      |          |          | main commit                       | conflicting_file |
      |          |          | Merge branch 'main' into feature1 |                  |
      | feature2 | local    | feature2 commit                   | feature2_file    |
      |          |          | main commit                       | conflicting_file |
      |          |          | Merge branch 'main' into feature2 |                  |
