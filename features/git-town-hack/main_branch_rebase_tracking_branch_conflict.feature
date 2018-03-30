Feature: git town-hack: resolving conflicts between main branch and its tracking branch

  As a developer creating a new feature branch while there are conflicting updates on the local and remote main branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given my repository has a feature branch named "existing-feature"
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
      |        | remote   | conflicting remote commit | conflicting_file | remote content |
    And I am on the "existing-feature" branch
    And my workspace has an uncommitted file
    When I run `git-town hack new-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH           | COMMAND                |
      | existing-feature | git fetch --prune      |
      |                  | git add -A             |
      |                  | git stash              |
      |                  | git checkout main      |
      | main             | git rebase origin/main |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And my repo has a rebase in progress
    And my uncommitted file is stashed


  Scenario: aborting
    When I run `git-town abort`
    Then it runs the commands
      | BRANCH           | COMMAND                       |
      | main             | git rebase --abort            |
      |                  | git checkout existing-feature |
      | existing-feature | git stash pop                 |
    And I end up on the "existing-feature" branch
    And my workspace has the uncommitted file again
    And there is no rebase in progress
    And my repository is left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run `git-town continue`
    Then it prints the error "You must resolve the conflicts before continuing"
    And my uncommitted file is stashed
    And my repo still has a rebase in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town continue`
    Then it runs the commands
      | BRANCH      | COMMAND                          |
      | main        | git rebase --continue            |
      |             | git push                         |
      |             | git checkout -b new-feature main |
      | new-feature | git stash pop                    |
    And I end up on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And now my repository has the following commits
      | BRANCH      | LOCATION         | MESSAGE                   | FILE NAME        |
      | main        | local and remote | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
      | new-feature | local            | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
    And now my repository has the following committed files
      | BRANCH      | NAME             | CONTENT          |
      | main        | conflicting_file | resolved content |
      | new-feature | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git-town continue`
    Then it runs the commands
      | BRANCH      | COMMAND                          |
      | main        | git push                         |
      |             | git checkout -b new-feature main |
      | new-feature | git stash pop                    |
    And I end up on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And now my repository has the following commits
      | BRANCH      | LOCATION         | MESSAGE                   | FILE NAME        |
      | main        | local and remote | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
      | new-feature | local            | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
    And now my repository has the following committed files
      | BRANCH      | NAME             | CONTENT          |
      | main        | conflicting_file | resolved content |
      | new-feature | conflicting_file | resolved content |
