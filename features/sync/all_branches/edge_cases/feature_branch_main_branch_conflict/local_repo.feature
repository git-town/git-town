Feature: handling merge conflicts between feature branch and main branch in a local repo

  Background:
    Given my repo does not have a remote origin
    And my repo has the local feature branches "feature-1", "feature-2", and "feature-3"
    And the following commits exist in my repo
      | BRANCH    | LOCATION | MESSAGE          | FILE NAME        | FILE CONTENT      |
      | main      | local    | main commit      | conflicting_file | main content      |
      | feature-1 | local    | feature-1 commit | feature1_file    | feature-1 content |
      | feature-2 | local    | feature-2 commit | conflicting_file | feature-2 content |
      | feature-3 | local    | feature-3 commit | feature3_file    | feature-3 content |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | main      | git add -A               |
      |           | git stash                |
      |           | git checkout feature-1   |
      | feature-1 | git merge --no-edit main |
      |           | git checkout feature-2   |
      | feature-2 | git merge --no-edit main |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      To continue by skipping the current branch, run "git-town skip".
      """
    And I am now on the "feature-2" branch
    And my uncommitted file is stashed
    And my repo now has a merge in progress

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH    | COMMAND                                       |
      | feature-2 | git merge --abort                             |
      |           | git checkout feature-1                        |
      | feature-1 | git reset --hard {{ sha 'feature-1 commit' }} |
      |           | git checkout main                             |
      | main      | git stash pop                                 |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo is left with my original commits

  Scenario: skipping
    When I run "git-town skip"
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | feature-2 | git merge --abort        |
      |           | git checkout feature-3   |
      | feature-3 | git merge --no-edit main |
      |           | git checkout main        |
      | main      | git stash pop            |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the following commits
      | BRANCH    | LOCATION | MESSAGE                            |
      | main      | local    | main commit                        |
      | feature-1 | local    | feature-1 commit                   |
      |           |          | main commit                        |
      |           |          | Merge branch 'main' into feature-1 |
      | feature-2 | local    | feature-2 commit                   |
      | feature-3 | local    | feature-3 commit                   |
      |           |          | main commit                        |
      |           |          | Merge branch 'main' into feature-3 |
    And my repo now has the following committed files
      | BRANCH    | NAME             | CONTENT           |
      | main      | conflicting_file | main content      |
      | feature-1 | conflicting_file | main content      |
      |           | feature1_file    | feature-1 content |
      | feature-2 | conflicting_file | feature-2 content |
      | feature-3 | conflicting_file | main content      |
      |           | feature3_file    | feature-3 content |

  Scenario: continuing without resolving the conflicts
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And I am still on the "feature-2" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress

  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | feature-2 | git commit --no-edit     |
      |           | git checkout feature-3   |
      | feature-3 | git merge --no-edit main |
      |           | git checkout main        |
      | main      | git stash pop            |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the following commits
      | BRANCH    | LOCATION | MESSAGE                            |
      | main      | local    | main commit                        |
      | feature-1 | local    | feature-1 commit                   |
      |           |          | main commit                        |
      |           |          | Merge branch 'main' into feature-1 |
      | feature-2 | local    | feature-2 commit                   |
      |           |          | main commit                        |
      |           |          | Merge branch 'main' into feature-2 |
      | feature-3 | local    | feature-3 commit                   |
      |           |          | main commit                        |
      |           |          | Merge branch 'main' into feature-3 |
    And my repo now has the following committed files
      | BRANCH    | NAME             | CONTENT           |
      | main      | conflicting_file | main content      |
      | feature-1 | conflicting_file | main content      |
      |           | feature1_file    | feature-1 content |
      | feature-2 | conflicting_file | resolved content  |
      | feature-3 | conflicting_file | main content      |
      |           | feature3_file    | feature-3 content |

  Scenario: continuing after resolving the conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    When I run "git-town continue"
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | feature-2 | git checkout feature-3   |
      | feature-3 | git merge --no-edit main |
      |           | git checkout main        |
      | main      | git stash pop            |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
