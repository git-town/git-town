Feature: Git Ship: handling conflicting remote feature branch updates when shipping the current feature branch


  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME          | FILE CONTENT               |
      | feature | remote   | remote conflicting commit | conflicting_file   | remote conflicting content |
      |         | local    | local conflicting commit  | conflicting_file   | local conflicting content  |
    And I run `git ship -m 'feature done'` while allowing errors


  Scenario: result
    Then I am still on the "feature" branch
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git ship --abort`
    Then I am still on the "feature" branch
    And there is no merge in progress
    And I still have the following commits
      | BRANCH  | LOCATION | MESSAGE                   | FILES              |
      | feature | local    | local conflicting commit  | conflicting_file   |
      |         | remote   | remote conflicting commit | conflicting_file   |
    And I still have the following committed files
      | BRANCH  | FILES              | CONTENT                   |
      | feature | conflicting_file   | local conflicting content |


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git ship --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                         |
      | feature | git commit --no-edit            |
      | feature | git merge --no-edit main        |
      | feature | git checkout main               |
      | main    | git merge --squash feature      |
      | main    | git commit -a -m 'feature done' |
      | main    | git push                        |
      | main    | git push origin :feature        |
      | main    | git branch -D feature           |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE      | FILES            |
      | main    | local and remote | feature done | conflicting_file |
    And now I have the following committed files
      | BRANCH  | FILES            |
      | main    | conflicting_file |


  Scenario: continuing after resolving conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git ship --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                         |
      | feature | git merge --no-edit main        |
      | feature | git checkout main               |
      | main    | git merge --squash feature      |
      | main    | git commit -a -m 'feature done' |
      | main    | git push                        |
      | main    | git push origin :feature        |
      | main    | git branch -D feature           |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE      | FILES            |
      | main    | local and remote | feature done | conflicting_file |
    And now I have the following committed files
      | BRANCH  | FILES            |
      | main    | conflicting_file |
