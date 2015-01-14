Feature: git extract: resolving conflicting remote main branch updates (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main    | remote   | conflicting remote commit | conflicting_file | remote content |
      |         | local    | conflicting local commit  | conflicting_file | local content  |
      | feature | local    | feature commit            | feature_file     |                |
      |         |          | refactor commit           | refactor_file    |                |
    And I am on the "feature" branch
    When I run `git extract refactor` with the last commit sha while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                |
      | feature | git fetch --prune      |
      | feature | git checkout main      |
      | main    | git rebase origin/main |
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git extract --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND              |
      | HEAD   | git rebase --abort   |
      | main   | git checkout feature |
    And I end up on the "feature" branch
    And there is no "refactor" branch
    And I am left with my original commits
    And there is no rebase in progress


  Scenario: continuing without resolving conflicts
    When I run `git extract --continue` while allowing errors
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git extract"
    And my repo has a rebase in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git extract --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND                               |
      | HEAD     | git rebase --continue                 |
      | main     | git push                              |
      | main     | git checkout -b refactor main         |
      | refactor | git cherry-pick [SHA:refactor commit] |
      | refactor | git push -u origin refactor           |
    And I end up on the "refactor" branch
    And now I have the following commits
      | BRANCH   | LOCATION         | MESSAGE                   | FILE NAME        |
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
      | BRANCH   | COMMAND                               |
      | main     | git push                              |
      | main     | git checkout -b refactor main         |
      | refactor | git cherry-pick [SHA:refactor commit] |
      | refactor | git push -u origin refactor           |
    And I end up on the "refactor" branch
    And now I have the following commits
      | BRANCH   | LOCATION         | MESSAGE                   | FILE NAME        |
      | main     | local and remote | conflicting remote commit | conflicting_file |
      |          |                  | conflicting local commit  | conflicting_file |
      | feature  | local            | feature commit            | feature_file     |
      |          |                  | refactor commit           | refactor_file    |
      | refactor | local and remote | conflicting remote commit | conflicting_file |
      |          |                  | conflicting local commit  | conflicting_file |
      |          |                  | refactor commit           | refactor_file    |
