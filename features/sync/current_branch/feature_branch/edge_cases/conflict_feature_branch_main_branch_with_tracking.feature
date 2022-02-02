Feature: handle conflicts between the current feature branch and the main branch (with tracking branch updates)

  Background:
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
      |         | remote   | feature commit             | feature_file     | feature content |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      To continue by skipping the current branch, run "git-town skip".
      """
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo now has a merge in progress

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH  | COMMAND                                                 |
      | feature | git merge --abort                                       |
      |         | git reset --hard {{ sha 'conflicting feature commit' }} |
      |         | git checkout main                                       |
      | main    | git checkout feature                                    |
      | feature | git stash pop                                           |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And there is no merge in progress
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local, remote | conflicting main commit    | conflicting_file | main content    |
      | feature | local         | conflicting feature commit | conflicting_file | feature content |
      |         | remote        | feature commit             | feature_file     | feature content |

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

  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git commit --no-edit |
      |         | git push             |
      |         | git stash pop        |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And all branches are now synchronized
    And my repo still has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
      |         | feature_file     | feature content  |

  Scenario: continuing after resolving the conflicts and comitting
    Given I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    When I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND       |
      | feature | git push      |
      |         | git stash pop |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And all branches are now synchronized
    And my repo still has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
      |         | feature_file     | feature content  |
