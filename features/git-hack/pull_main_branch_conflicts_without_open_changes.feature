Feature: git-hack handling conflicting remote main branch updates with open changes

  Background:
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | remote   | conflicting remote commit | conflicting_file | remote content |
      |        | local    | conflicting local commit  | conflicting_file | local content  |
    And I am on the "feature" branch
    When I run `git hack other_feature` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                |
      | feature | git checkout main      |
      | main    | git fetch --prune      |
      | main    | git rebase origin/main |
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git hack --abort`
    Then it runs the Git commands
      | BRANCH  | COMMAND              |
      | HEAD    | git rebase --abort   |
      | main    | git checkout feature |
    And I end up on the "feature" branch
    And there is no rebase in progress
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE                   | FILES            |
      | main   | remote   | conflicting remote commit | conflicting_file |
      |        | local    | conflicting local commit  | conflicting_file |


  Scenario: continuing without resolving conflicts
    When I run `git hack --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git hack"
    And my repo still has a rebase in progress


  Scenario Outline: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `<command>`
    Then I end up on the "other_feature" branch
    And now I have the following commits
      | BRANCH        | LOCATION         | MESSAGE                   | FILES            |
      | main          | local and remote | conflicting remote commit | conflicting_file |
      |               |                  | conflicting local commit  | conflicting_file |
      | other_feature | local            | conflicting remote commit | conflicting_file |
      |               |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH        | FILES            | CONTENT          |
      | main          | conflicting_file | resolved content |
      | other_feature | conflicting_file | resolved content |

    Examples:
      | command                                    |
      | git hack --continue                        |
      | git rebase --continue; git hack --continue |
