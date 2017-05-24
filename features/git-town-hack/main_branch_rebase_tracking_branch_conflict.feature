Feature: git town-hack: resolving conflicts between main branch and its tracking branch

  As a developer creating a new feature branch while there are conflicting updates on the local and remote main branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given I have a feature branch named "existing-feature"
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
      |        | remote   | conflicting remote commit | conflicting_file | remote content |
    And I am on the "existing-feature" branch
    And I have an uncommitted file
    When I run `git-town hack new-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH           | COMMAND                |
      | existing-feature | git fetch --prune      |
      |                  | git add -A             |
      |                  | git stash              |
      |                  | git checkout main      |
      | main             | git rebase origin/main |
    And I get the error
      """
      To abort, run "git-town hack --abort".
      To continue after you have resolved the conflicts, run "git-town hack --continue".
      """
    And my repo has a rebase in progress
    And my uncommitted file is stashed


  Scenario: aborting
    When I run `git-town hack --abort`
    Then it runs the commands
      | BRANCH           | COMMAND                       |
      | main             | git rebase --abort            |
      |                  | git checkout existing-feature |
      | existing-feature | git stash pop                 |
    And I end up on the "existing-feature" branch
    And I again have my uncommitted file
    And there is no rebase in progress
    And I am left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run `git-town hack --continue`
    Then I get the error "You must resolve the conflicts before continuing"
    And my uncommitted file is stashed
    And my repo still has a rebase in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town hack --continue`
    Then it runs the commands
      | BRANCH      | COMMAND                          |
      | main        | git rebase --continue            |
      |             | git push                         |
      |             | git checkout -b new-feature main |
      | new-feature | git push -u origin new-feature   |
      |             | git stash pop                    |
    And I end up on the "new-feature" branch
    And I still have my uncommitted file
    And now I have the following commits
      | BRANCH      | LOCATION         | MESSAGE                   | FILE NAME        |
      | main        | local and remote | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
      | new-feature | local and remote | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH      | NAME             | CONTENT          |
      | main        | conflicting_file | resolved content |
      | new-feature | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git-town hack --continue`
    Then it runs the commands
      | BRANCH      | COMMAND                          |
      | main        | git push                         |
      |             | git checkout -b new-feature main |
      | new-feature | git push -u origin new-feature   |
      |             | git stash pop                    |
    And I end up on the "new-feature" branch
    And I still have my uncommitted file
    And now I have the following commits
      | BRANCH      | LOCATION         | MESSAGE                   | FILE NAME        |
      | main        | local and remote | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
      | new-feature | local and remote | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH      | NAME             | CONTENT          |
      | main        | conflicting_file | resolved content |
      | new-feature | conflicting_file | resolved content |
