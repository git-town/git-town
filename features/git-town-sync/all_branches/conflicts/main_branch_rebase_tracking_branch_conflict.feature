Feature: git-town sync --all: handling rebase conflicts between main branch and its tracking branch

  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        | FILE CONTENT        |
      | main    | local    | main local commit  | conflicting_file | main local content  |
      | main    | remote   | main remote commit | conflicting_file | main remote content |
      | feature | local    | feature commit     | feature_file     | feature content     |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run `git-town sync --all`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                |
      | main   | git fetch --prune      |
      |        | git add -A             |
      |        | git stash              |
      |        | git rebase origin/main |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And my uncommitted file is stashed
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git-town abort`
    Then it runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git stash pop      |
    And I end up on the "main" branch
    And my workspace has the uncommitted file again
    And my repository has the following commits
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | local    | main local commit  | conflicting_file |
      |         | remote   | main remote commit | conflicting_file |
      | feature | local    | feature commit     | feature_file     |


  Scenario: continuing without resolving the conflicts
    When I run `git-town continue`
    Then it runs no commands
    And it prints the error "You must resolve the conflicts before continuing"
    And my uncommitted file is stashed
    And my repo still has a rebase in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git-town continue`
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git rebase --continue              |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git checkout main                  |
      | main    | git push --tags                    |
      |         | git stash pop                      |
    And I end up on the "main" branch
    And my workspace has the uncommitted file again
    And my repository has the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME        |
      | main    | local and remote | main remote commit               | conflicting_file |
      |         |                  | main local commit                | conflicting_file |
      | feature | local and remote | feature commit                   | feature_file     |
      |         |                  | main remote commit               | conflicting_file |
      |         |                  | main local commit                | conflicting_file |
      |         |                  | Merge branch 'main' into feature |                  |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run `git rebase --continue; git-town continue`
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git checkout main                  |
      | main    | git push --tags                    |
      |         | git stash pop                      |
    And I end up on the "main" branch
    And my workspace has the uncommitted file again
    And my repository has the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME        |
      | main    | local and remote | main remote commit               | conflicting_file |
      |         |                  | main local commit                | conflicting_file |
      | feature | local and remote | feature commit                   | feature_file     |
      |         |                  | main remote commit               | conflicting_file |
      |         |                  | main local commit                | conflicting_file |
      |         |                  | Merge branch 'main' into feature |                  |
