Feature: Syncing before creating the pull request

  Background:
    Given my repo has a local feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME        | FILE CONTENT    |
      | main    | local, remote | main commit    | conflicting_file | main_content    |
      | feature | local         | feature commit | conflicting_file | feature content |
    And my computer has the "open" tool installed
    And my repo's origin is "git@github.com:git-town/git-town.git"
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town new-pull-request"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git add -A               |
      |         | git stash                |
      |         | git checkout main        |
      | main    | git rebase origin/main   |
      |         | git checkout feature     |
      | feature | git merge --no-edit main |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo now has a merge in progress

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git merge --abort    |
      |         | git checkout main    |
      | main    | git checkout feature |
      | feature | git stash pop        |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And there is no merge in progress
    And my repo is left with my original commits

  Scenario: continuing without resolving the conflicts
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress

  @skipWindows
  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git commit --no-edit                                               |
      |         | git push -u origin feature                                         |
      |         | git stash pop                                                      |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                          | FILE NAME        |
      | main    | local, remote | main commit                      | conflicting_file |
      | feature | local, remote | feature commit                   | conflicting_file |
      |         |               | main commit                      | conflicting_file |
      |         |               | Merge branch 'main' into feature |                  |

  @skipWindows
  Scenario: continuing after resolving conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    When I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git push -u origin feature                                         |
      |         | git stash pop                                                      |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                          | FILE NAME        |
      | main    | local, remote | main commit                      | conflicting_file |
      | feature | local, remote | feature commit                   | conflicting_file |
      |         |               | main commit                      | conflicting_file |
      |         |               | Merge branch 'main' into feature |                  |
