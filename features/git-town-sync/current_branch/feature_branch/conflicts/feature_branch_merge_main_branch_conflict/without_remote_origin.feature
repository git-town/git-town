Feature: git-town sync: resolving conflicts between the current feature branch and the main branch (without remote repo)

  Background:
    Given my repo does not have a remote origin
    And my repository has a local feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run `git-town sync`


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git add -A               |
      |         | git stash                |
      |         | git merge --no-edit main |
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
      | BRANCH  | COMMAND           |
      | feature | git merge --abort |
      |         | git stash pop     |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And there is no merge in progress
    And my repository is left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run `git-town continue`
    Then it runs no commands
    And it prints the error "You must resolve the conflicts before continuing"
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town continue`
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git commit --no-edit |
      |         | git stash pop        |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And my repository still has the following commits
      | BRANCH  | LOCATION | MESSAGE                          | FILE NAME        |
      | main    | local    | conflicting main commit          | conflicting_file |
      | feature | local    | conflicting feature commit       | conflicting_file |
      |         |          | conflicting main commit          | conflicting_file |
      |         |          | Merge branch 'main' into feature |                  |
    And my repository still has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts and comitting
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git-town continue`
    Then it runs the commands
      | BRANCH  | COMMAND       |
      | feature | git stash pop |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And my repository still has the following commits
      | BRANCH  | LOCATION | MESSAGE                          | FILE NAME        |
      | main    | local    | conflicting main commit          | conflicting_file |
      | feature | local    | conflicting feature commit       | conflicting_file |
      |         |          | conflicting main commit          | conflicting_file |
      |         |          | Merge branch 'main' into feature |                  |
    And my repository still has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
