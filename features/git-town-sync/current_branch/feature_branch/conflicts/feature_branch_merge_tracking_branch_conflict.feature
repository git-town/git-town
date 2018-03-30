Feature: git-town sync: resolving conflicts between the current feature branch and its tracking branch

  As a developer syncing a feature branch that conflicts with the tracking branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | local    | local conflicting commit  | conflicting_file | local conflicting content  |
      |         | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run `git-town sync`


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      To continue by skipping the current branch, run "git-town skip".
      """
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git-town abort`
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git merge --abort    |
      |         | git checkout main    |
      | main    | git checkout feature |
      | feature | git stash pop        |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And there is no merge in progress
    And my repository is left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run `git-town continue`
    Then it prints the error "You must resolve the conflicts before continuing"
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town continue`
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git commit --no-edit     |
      |         | git merge --no-edit main |
      |         | git push                 |
      |         | git stash pop            |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And now my repository has the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME        |
      | feature | local and remote | local conflicting commit                                   | conflicting_file |
      |         |                  | remote conflicting commit                                  | conflicting_file |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                  |
    And now my repository has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | feature | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git-town continue`
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git merge --no-edit main |
      |         | git push                 |
      |         | git stash pop            |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And now my repository has the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME        |
      | feature | local and remote | local conflicting commit                                   | conflicting_file |
      |         |                  | remote conflicting commit                                  | conflicting_file |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                  |
    And now my repository has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | feature | conflicting_file | resolved content |
