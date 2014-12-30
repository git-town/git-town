Feature: git-sync-all: handling rebase conflicts between non-feature branch and its tracking branch without open changes

  Background:
    Given I have branches named "qa" and "production"
    And my non-feature branches are "qa" and "production"
    And the following commits exist in my repository
      | branch     | location         | message                  | file name        | file content              |
      | main       | remote           | main commit              | main_file        | main content              |
      | production | local            | production local commit  | conflicting_file | production local content  |
      |            | remote           | production remote commit | conflicting_file | production remote content |
      | qa         | local and remote | qa commit                | qa_file          | qa content                |
    And I am on the "main" branch
    When I run `git sync-all` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH     | COMMAND                      |
      | main       | git fetch --prune            |
      | main       | git rebase origin/main       |
      | main       | git checkout production      |
      | production | git rebase origin/production |
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git sync-all --abort`
    Then it runs the Git commands
      | BRANCH     | COMMAND            |
      | HEAD       | git rebase --abort |
      | production | git checkout main  |
    And I end up on the "main" branch
    And I have the following commits
      | branch     | location         | message                  | files            |
      | main       | local and remote | main commit              | main_file        |
      | production | local            | production local commit  | conflicting_file |
      |            | remote           | production remote commit | conflicting_file |
      | qa         | local and remote | qa commit                | qa_file          |


  Scenario: skipping
    When I run `git sync-all --skip`
    Then it runs the Git commands
      | BRANCH     | COMMAND              |
      | HEAD       | git rebase --abort   |
      | production | git checkout qa      |
      | qa         | git rebase origin/qa |
      | qa         | git checkout main    |
    And I end up on the "main" branch
    And I have the following commits
      | branch     | location         | message                  | files            |
      | main       | local and remote | main commit              | main_file        |
      | production | local            | production local commit  | conflicting_file |
      |            | remote           | production remote commit | conflicting_file |
      | qa         | local and remote | qa commit                | qa_file          |


  Scenario: continuing without resolving conflicts
    When I run `git sync-all --continue` while allowing errors
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And my repo still has a rebase in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync-all --continue`
    Then it runs the Git commands
      | BRANCH     | COMMAND               |
      | HEAD       | git rebase --continue |
      | production | git push              |
      | production | git checkout qa       |
      | qa         | git rebase origin/qa  |
      | qa         | git checkout main     |
    And I end up on the "main" branch
    And I have the following commits
      | branch     | location         | message                  | files            |
      | main       | local and remote | main commit              | main_file        |
      | production | local and remote | production local commit  | conflicting_file |
      |            | local and remote | production remote commit | conflicting_file |
      | qa         | local and remote | qa commit                | qa_file          |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run `git rebase --continue; git sync-all --continue`
    Then it runs the Git commands
      | BRANCH     | COMMAND              |
      | production | git push             |
      | production | git checkout qa      |
      | qa         | git rebase origin/qa |
      | qa         | git checkout main    |
    And I end up on the "main" branch
    And I have the following commits
      | branch     | location         | message                  | files            |
      | main       | local and remote | main commit              | main_file        |
      | production | local and remote | production local commit  | conflicting_file |
      |            | local and remote | production remote commit | conflicting_file |
      | qa         | local and remote | qa commit                | qa_file          |
