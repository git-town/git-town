Feature: Syncing before creating the pull request

  As a developer trying to create a pull request when my feature branch conflicts with the main branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME        | FILE CONTENT    |
      | main    | local and remote | main commit    | conflicting_file | main_content    |
      | feature | local            | feature commit | conflicting_file | feature content |
    And I have "open" installed
    And my repo's remote origin is git@github.com:Originate/git-town.git
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run `git-town new-pull-request`


  @finishes-with-non-empty-stash
  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And it prints the error:
      """
      To abort, run "git-town new-pull-request --abort".
      To continue after you have resolved the conflicts, run "git-town new-pull-request --continue".
      """
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git-town new-pull-request --abort`
    Then Git Town runs the commands
      | BRANCH  | COMMAND              |
      | feature | git merge --abort    |
      |         | git checkout main    |
      | main    | git checkout feature |
      | feature | git stash pop        |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And there is no merge in progress
    And my repository is left with my original commits


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving the conflicts
    When I run `git-town new-pull-request --continue`
    Then Git Town runs no commands
    And it prints the error "You must resolve the conflicts before continuing"
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town new-pull-request --continue`
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                             |
      | feature | git commit --no-edit                                                |
      |         | git push                                                            |
      |         | git stash pop                                                       |
      | <none>  | open https://github.com/Originate/git-town/compare/feature?expand=1 |
    And I see a new GitHub pull request for the "feature" branch in the "Originate/git-town" repo in my browser
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME        |
      | main    | local and remote | main commit                      | conflicting_file |
      | feature | local and remote | feature commit                   | conflicting_file |
      |         |                  | main commit                      | conflicting_file |
      |         |                  | Merge branch 'main' into feature |                  |


  Scenario: continuing after resolving conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git-town new-pull-request --continue`
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                             |
      | feature | git push                                                            |
      |         | git stash pop                                                       |
      | <none>  | open https://github.com/Originate/git-town/compare/feature?expand=1 |
    And I see a new GitHub pull request for the "feature" branch in the "Originate/git-town" repo in my browser
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME        |
      | main    | local and remote | main commit                      | conflicting_file |
      | feature | local and remote | feature commit                   | conflicting_file |
      |         |                  | main commit                      | conflicting_file |
      |         |                  | Merge branch 'main' into feature |                  |
