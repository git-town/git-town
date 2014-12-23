Feature: Git Ship: handling merge conflicts between feature and main branch when shipping the current feature branch

  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I run `git ship -m 'feature done'` while allowing errors


  Scenario: result
    Then I am still on the "feature" branch
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git ship --abort`
    Then I am still on the "feature" branch
    And there is no merge in progress
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                    | FILES            |
      | main    | local and remote | conflicting main commit    | conflicting_file |
      | feature | local            | conflicting feature commit | conflicting_file |
    And I still have the following committed files
      | BRANCH  | FILES            | CONTENT         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git ship --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                         |
      | feature | git commit --no-edit            |
      | feature | git checkout main               |
      | main    | git merge --squash feature      |
      | main    | git commit -a -m 'feature done' |
      | main    | git push                        |
      | main    | git push origin :feature        |
      | main    | git branch -D feature           |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                 | FILES            |
      | main    | local and remote | conflicting main commit | conflicting_file |
      |         |                  | feature done            | conflicting_file |
    And now I have the following committed files
      | BRANCH  | FILES            |
      | main    | conflicting_file |


  Scenario: continuing after resolving conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git ship --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                         |
      | feature | git checkout main               |
      | main    | git merge --squash feature      |
      | main    | git commit -a -m 'feature done' |
      | main    | git push                        |
      | main    | git push origin :feature        |
      | main    | git branch -D feature           |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                 | FILES            |
      | main    | local and remote | conflicting main commit | conflicting_file |
      |         |                  | feature done            | conflicting_file |
    And now I have the following committed files
      | BRANCH  | FILES            |
      | main    | conflicting_file |
