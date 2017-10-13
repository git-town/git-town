Feature: git-town sync --all: handling rebase conflicts between perennial branch and its tracking branch

  Background:
    Given my repository has the perennial branches "production" and "qa"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE                  | FILE NAME        | FILE CONTENT              |
      | main       | remote           | main commit              | main_file        | main content              |
      | production | local            | production local commit  | conflicting_file | production local content  |
      |            | remote           | production remote commit | conflicting_file | production remote content |
      | qa         | local and remote | qa commit                | qa_file          | qa content                |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run `git-town sync --all`


  Scenario: result
    Then I am not prompted for any parent branches
    And Git Town runs the commands
      | BRANCH     | COMMAND                      |
      | main       | git fetch --prune            |
      |            | git add -A                   |
      |            | git stash                    |
      |            | git rebase origin/main       |
      |            | git checkout production      |
      | production | git rebase origin/production |
    And it prints the error:
      """
      To abort, run "git-town sync --abort".
      To continue after you have resolved the conflicts, run "git-town sync --continue".
      To skip the sync of the 'production' branch, run "git-town sync --skip".
      """
    And my uncommitted file is stashed
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git-town sync --abort`
    Then Git Town runs the commands
      | BRANCH     | COMMAND            |
      | production | git rebase --abort |
      |            | git checkout main  |
      | main       | git stash pop      |
    And I end up on the "main" branch
    And my workspace has the uncommitted file again
    And my repository has the following commits
      | BRANCH     | LOCATION         | MESSAGE                  | FILE NAME        |
      | main       | local and remote | main commit              | main_file        |
      | production | local            | production local commit  | conflicting_file |
      |            | remote           | production remote commit | conflicting_file |
      | qa         | local and remote | qa commit                | qa_file          |


  Scenario: skipping
    When I run `git-town sync --skip`
    Then Git Town runs the commands
      | BRANCH     | COMMAND              |
      | production | git rebase --abort   |
      |            | git checkout qa      |
      | qa         | git rebase origin/qa |
      |            | git checkout main    |
      | main       | git push --tags      |
      |            | git stash pop        |
    And I end up on the "main" branch
    And my workspace has the uncommitted file again
    And my repository has the following commits
      | BRANCH     | LOCATION         | MESSAGE                  | FILE NAME        |
      | main       | local and remote | main commit              | main_file        |
      | production | local            | production local commit  | conflicting_file |
      |            | remote           | production remote commit | conflicting_file |
      | qa         | local and remote | qa commit                | qa_file          |


  Scenario: continuing without resolving the conflicts
    When I run `git-town sync --continue`
    Then Git Town runs no commands
    And it prints the error "You must resolve the conflicts before continuing"
    And my uncommitted file is stashed
    And my repo still has a rebase in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git-town sync --continue`
    Then Git Town runs the commands
      | BRANCH     | COMMAND               |
      | production | git rebase --continue |
      |            | git push              |
      |            | git checkout qa       |
      | qa         | git rebase origin/qa  |
      |            | git checkout main     |
      | main       | git push --tags       |
      |            | git stash pop         |
    And I end up on the "main" branch
    And my workspace has the uncommitted file again
    And my repository has the following commits
      | BRANCH     | LOCATION         | MESSAGE                  | FILE NAME        |
      | main       | local and remote | main commit              | main_file        |
      | production | local and remote | production remote commit | conflicting_file |
      |            |                  | production local commit  | conflicting_file |
      | qa         | local and remote | qa commit                | qa_file          |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run `git rebase --continue; git-town sync --continue`
    Then Git Town runs the commands
      | BRANCH     | COMMAND              |
      | production | git push             |
      |            | git checkout qa      |
      | qa         | git rebase origin/qa |
      |            | git checkout main    |
      | main       | git push --tags      |
      |            | git stash pop        |
    And I end up on the "main" branch
    And my workspace has the uncommitted file again
    And my repository has the following commits
      | BRANCH     | LOCATION         | MESSAGE                  | FILE NAME        |
      | main       | local and remote | main commit              | main_file        |
      | production | local and remote | production remote commit | conflicting_file |
      |            |                  | production local commit  | conflicting_file |
      | qa         | local and remote | qa commit                | qa_file          |
