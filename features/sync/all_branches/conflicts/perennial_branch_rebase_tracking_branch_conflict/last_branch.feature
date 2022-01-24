Feature: git-town sync --all: handling rebase conflicts between perennial branch and its tracking branch

  Background:
    Given my repo has the perennial branches "production" and "qa"
    And the following commits exist in my repo
      | BRANCH     | LOCATION      | MESSAGE           | FILE NAME        | FILE CONTENT       |
      | main       | remote        | main commit       | main_file        | main content       |
      | production | local, remote | production commit | production_file  | production content |
      | qa         | local         | qa local commit   | conflicting_file | qa local content   |
      |            | remote        | qa remote commit  | conflicting_file | qa remote content  |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then I am not prompted for any parent branches
    And it runs the commands
      | BRANCH     | COMMAND                      |
      | main       | git fetch --prune --tags     |
      |            | git add -A                   |
      |            | git stash                    |
      |            | git rebase origin/main       |
      |            | git checkout production      |
      | production | git rebase origin/production |
      |            | git checkout qa              |
      | qa         | git rebase origin/qa         |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      To continue by skipping the current branch, run "git-town skip".
      """
    And my uncommitted file is stashed
    And my repo now has a rebase in progress

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH     | COMMAND                 |
      | qa         | git rebase --abort      |
      |            | git checkout production |
      | production | git checkout main       |
      | main       | git stash pop           |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE           | FILE NAME        |
      | main       | local, remote | main commit       | main_file        |
      | production | local, remote | production commit | production_file  |
      | qa         | local         | qa local commit   | conflicting_file |
      |            | remote        | qa remote commit  | conflicting_file |

  Scenario: skipping
    When I run "git-town skip"
    Then it runs the commands
      | BRANCH | COMMAND            |
      | qa     | git rebase --abort |
      |        | git checkout main  |
      | main   | git push --tags    |
      |        | git stash pop      |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE           | FILE NAME        |
      | main       | local, remote | main commit       | main_file        |
      | production | local, remote | production commit | production_file  |
      | qa         | local         | qa local commit   | conflicting_file |
      |            | remote        | qa remote commit  | conflicting_file |

  Scenario: continuing without resolving the conflicts
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And my uncommitted file is stashed
    And my repo still has a rebase in progress

  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND               |
      | qa     | git rebase --continue |
      |        | git push              |
      |        | git checkout main     |
      | main   | git push --tags       |
      |        | git stash pop         |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE           | FILE NAME        |
      | main       | local, remote | main commit       | main_file        |
      | production | local, remote | production commit | production_file  |
      | qa         | local, remote | qa remote commit  | conflicting_file |
      |            |               | qa local commit   | conflicting_file |

  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | qa     | git push          |
      |        | git checkout main |
      | main   | git push --tags   |
      |        | git stash pop     |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE           | FILE NAME        |
      | main       | local, remote | main commit       | main_file        |
      | production | local, remote | production commit | production_file  |
      | qa         | local, remote | qa remote commit  | conflicting_file |
      |            |               | qa local commit   | conflicting_file |
