Feature: git hack: resolving conflicts between main branch and its tracking branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have a feature branch named "existing_feature"
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | remote   | conflicting remote commit | conflicting_file | remote content |
      |        | local    | conflicting local commit  | conflicting_file | local content  |
    And I am on the "existing_feature" branch
    When I run `git hack new_feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH           | COMMAND                |
      | existing_feature | git fetch --prune      |
      | existing_feature | git checkout main      |
      | main             | git rebase origin/main |
    And I get the error
      """
      To abort, run "git hack --abort".
      To continue after you have resolved the conflicts, run "git hack --continue".
      """
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git hack --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND                       |
      | main   | git rebase --abort            |
      | main   | git checkout existing_feature |
    And I end up on the "existing_feature" branch
    And there is no rebase in progress
    And I am left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run `git hack --continue`
    Then I get the error "You must resolve the conflicts before continuing the git hack"
    And my repo still has a rebase in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git hack --continue `
    Then it runs the Git commands
      | BRANCH | COMMAND                          |
      | main   | git rebase --continue            |
      | main   | git push                         |
      | main   | git checkout -b new_feature main |
    And I end up on the "new_feature" branch
    And now I have the following commits
      | BRANCH      | LOCATION         | MESSAGE                   | FILE NAME        |
      | main        | local and remote | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
      | new_feature | local            | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH      | FILES            | CONTENT          |
      | main        | conflicting_file | resolved content |
      | new_feature | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git hack --continue `
    Then it runs the Git commands
      | BRANCH | COMMAND                          |
      | main   | git push                         |
      | main   | git checkout -b new_feature main |
    And I end up on the "new_feature" branch
    And now I have the following commits
      | BRANCH      | LOCATION         | MESSAGE                   | FILE NAME        |
      | main        | local and remote | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
      | new_feature | local            | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH      | FILES            | CONTENT          |
      | main        | conflicting_file | resolved content |
      | new_feature | conflicting_file | resolved content |
