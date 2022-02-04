Feature: handle rebase conflicts between perennial branch and its tracking branch

  Background:
    Given my repo has the perennial branches "perennial-1", "perennial-2", and "perennial-3"
    And my repo contains the commits
      | BRANCH      | LOCATION      | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main        | remote        | main commit               | main_file        | main content               |
      | perennial-1 | local, remote | perennial-1 commit        | peren1_file      | perennial-1 content        |
      | perennial-2 | local         | perennial-2 local commit  | conflicting_file | perennial-2 local content  |
      |             | remote        | perennial-2 remote commit | conflicting_file | perennial-2 remote content |
      | perennial-3 | local, remote | perennial-3 commit        | peren3_file      | perennial-3 content        |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then I am not prompted for any parent branches
    And it runs the commands
      | BRANCH      | COMMAND                       |
      | main        | git fetch --prune --tags      |
      |             | git add -A                    |
      |             | git stash                     |
      |             | git rebase origin/main        |
      |             | git checkout perennial-1      |
      | perennial-1 | git rebase origin/perennial-1 |
      |             | git checkout perennial-2      |
      | perennial-2 | git rebase origin/perennial-2 |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      To continue by skipping the current branch, run "git-town skip".
      """
    And my uncommitted file is stashed
    And my repo now has a rebase in progress

  Scenario: abort
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH      | COMMAND                  |
      | perennial-2 | git rebase --abort       |
      |             | git checkout perennial-1 |
      | perennial-1 | git checkout main        |
      | main        | git stash pop            |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the commits
      | BRANCH      | LOCATION      | MESSAGE                   |
      | main        | local, remote | main commit               |
      | perennial-1 | local, remote | perennial-1 commit        |
      | perennial-2 | local         | perennial-2 local commit  |
      |             | remote        | perennial-2 remote commit |
      | perennial-3 | local, remote | perennial-3 commit        |

  Scenario: skip
    When I run "git-town skip"
    Then it runs the commands
      | BRANCH      | COMMAND                       |
      | perennial-2 | git rebase --abort            |
      |             | git checkout perennial-3      |
      | perennial-3 | git rebase origin/perennial-3 |
      |             | git checkout main             |
      | main        | git push --tags               |
      |             | git stash pop                 |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the commits
      | BRANCH      | LOCATION      | MESSAGE                   |
      | main        | local, remote | main commit               |
      | perennial-1 | local, remote | perennial-1 commit        |
      | perennial-2 | local         | perennial-2 local commit  |
      |             | remote        | perennial-2 remote commit |
      | perennial-3 | local, remote | perennial-3 commit        |

  Scenario: continue without resolving the conflicts
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And my uncommitted file is stashed
    And my repo still has a rebase in progress

  Scenario: continue after resolving the conflicts
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH      | COMMAND                       |
      | perennial-2 | git rebase --continue         |
      |             | git push                      |
      |             | git checkout perennial-3      |
      | perennial-3 | git rebase origin/perennial-3 |
      |             | git checkout main             |
      | main        | git push --tags               |
      |             | git stash pop                 |
    And all branches are now synchronized
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And there is no rebase in progress anymore

  Scenario: continue after resolving the conflicts and continuing the rebase
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH      | COMMAND                       |
      | perennial-2 | git push                      |
      |             | git checkout perennial-3      |
      | perennial-3 | git rebase origin/perennial-3 |
      |             | git checkout main             |
      | main        | git push --tags               |
      |             | git stash pop                 |
