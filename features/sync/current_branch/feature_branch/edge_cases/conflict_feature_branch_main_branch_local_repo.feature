Feature: handle conflicts between the current feature branch and the main branch (without remote repo)

  Background:
    Given my repo does not have a remote origin
    And my repo has a local feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

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
    And my repo now has a merge in progress

  Scenario: abort
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH  | COMMAND           |
      | feature | git merge --abort |
      |         | git stash pop     |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And there is no merge in progress
    And my repo is left with my original commits

  Scenario: continue without resolving the conflicts
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress

  Scenario: continue after resolving the conflicts
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git commit --no-edit |
      |         | git stash pop        |
    And all branches are now synchronized
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And my repo now has these committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

  Scenario: continue after resolving the conflicts and comitting
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND       |
      | feature | git stash pop |
