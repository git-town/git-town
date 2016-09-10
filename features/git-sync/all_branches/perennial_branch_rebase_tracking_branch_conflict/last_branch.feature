Feature: git sync --all: handling rebase conflicts between perennial branch and its tracking branch

  Background:
    Given I have perennial branches named "production" and "qa"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE           | FILE NAME        | FILE CONTENT       |
      | main       | remote           | main commit       | main_file        | main content       |
      | production | local and remote | production commit | production_file  | production content |
      | qa         | local            | qa local commit   | conflicting_file | qa local content   |
      |            | remote           | qa remote commit  | conflicting_file | qa remote content  |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git sync --all`


  Scenario: result
    Then I am not prompted for any parent branches
    And it runs the commands
      | BRANCH     | COMMAND                      |
      | main       | git fetch --prune            |
      |            | git stash -a                 |
      |            | git rebase origin/main       |
      |            | git checkout production      |
      | production | git rebase origin/production |
      |            | git checkout qa              |
      | qa         | git rebase origin/qa         |
    And I get the error
      """
      To abort, run "git sync --abort".
      To continue after you have resolved the conflicts, run "git sync --continue".
      To skip the sync of the 'qa' branch, run "git sync --skip".
      """
    And my uncommitted file is stashed
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the commands
      | BRANCH     | COMMAND                 |
      | qa         | git rebase --abort      |
      |            | git checkout production |
      | production | git checkout main       |
      | main       | git stash pop           |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE           | FILE NAME        |
      | main       | local and remote | main commit       | main_file        |
      | production | local and remote | production commit | production_file  |
      | qa         | local            | qa local commit   | conflicting_file |
      |            | remote           | qa remote commit  | conflicting_file |


  Scenario: skipping
    When I run `git sync --skip`
    Then it runs the commands
      | BRANCH | COMMAND            |
      | qa     | git rebase --abort |
      |        | git checkout main  |
      | main   | git push --tags    |
      |        | git stash pop      |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE           | FILE NAME        |
      | main       | local and remote | main commit       | main_file        |
      | production | local and remote | production commit | production_file  |
      | qa         | local            | qa local commit   | conflicting_file |
      |            | remote           | qa remote commit  | conflicting_file |


  Scenario: continuing without resolving the conflicts
    When I run `git sync --continue`
    Then it runs no commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And my uncommitted file is stashed
    And my repo still has a rebase in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then it runs the commands
      | BRANCH | COMMAND               |
      | qa     | git rebase --continue |
      |        | git push              |
      |        | git checkout main     |
      | main   | git push --tags       |
      |        | git stash pop         |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE           | FILE NAME        |
      | main       | local and remote | main commit       | main_file        |
      | production | local and remote | production commit | production_file  |
      | qa         | local and remote | qa remote commit  | conflicting_file |
      |            |                  | qa local commit   | conflicting_file |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run `git rebase --continue; git sync --continue`
    Then it runs the commands
      | BRANCH | COMMAND           |
      | qa     | git push          |
      |        | git checkout main |
      | main   | git push --tags   |
      |        | git stash pop     |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE           | FILE NAME        |
      | main       | local and remote | main commit       | main_file        |
      | production | local and remote | production commit | production_file  |
      | qa         | local and remote | qa remote commit  | conflicting_file |
      |            |                  | qa local commit   | conflicting_file |
