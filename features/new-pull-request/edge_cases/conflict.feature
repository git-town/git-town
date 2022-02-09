Feature: merge conflict

  Background:
    Given a local feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME        | FILE CONTENT    |
      | main    | local, origin | main commit    | conflicting_file | main content    |
      | feature | local         | feature commit | conflicting_file | feature content |
    And my computer has the "open" tool installed
    And my repo's origin is "git@github.com:git-town/git-town.git"
    And I am on the "feature" branch
    When I run "git-town new-pull-request"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
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
    And my repo now has a merge in progress

  Scenario: abort
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git merge --abort    |
      |         | git checkout main    |
      | main    | git checkout feature |
    And I am still on the "feature" branch
    And there is no merge in progress
    And now the initial commits exist

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And I am still on the "feature" branch
    And my repo still has a merge in progress

  @skipWindows
  Scenario: resolve and continue
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git commit --no-edit                                               |
      |         | git push -u origin feature                                         |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """
    And I am still on the "feature" branch
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | main commit                      |
      | feature | local, origin | feature commit                   |
      |         |               | main commit                      |
      |         |               | Merge branch 'main' into feature |
    And my repo now has these committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

  @skipWindows
  Scenario: resolve, commit, and continue
    Given I resolve the conflict in "conflicting_file"
    When I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git push -u origin feature                                         |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """
    And I am still on the "feature" branch
