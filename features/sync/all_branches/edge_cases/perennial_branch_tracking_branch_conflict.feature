Feature: handle rebase conflicts between perennial branch and its tracking branch

  Background:
    Given my repo has the perennial branches "alpha", "beta", and "gamma"
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT        |
      | main   | remote        | main commit        | main_file        | main content        |
      | alpha  | local, remote | alpha commit       | alpha_file       | alpha content       |
      | beta   | local         | local beta commit  | conflicting_file | local beta content  |
      |        | remote        | remote beta commit | conflicting_file | remote beta content |
      | gamma  | local, remote | gamma commit       | gamma_file       | gamma content       |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then I am not prompted for any parent branches
    And it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/main   |
      |        | git checkout alpha       |
      | alpha  | git rebase origin/alpha  |
      |        | git checkout beta        |
      | beta   | git rebase origin/beta   |
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
      | BRANCH | COMMAND            |
      | beta   | git rebase --abort |
      |        | git checkout alpha |
      | alpha  | git checkout main  |
      | main   | git stash pop      |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE            |
      | main   | local, remote | main commit        |
      | alpha  | local, remote | alpha commit       |
      | beta   | local         | local beta commit  |
      |        | remote        | remote beta commit |
      | gamma  | local, remote | gamma commit       |

  Scenario: skip
    When I run "git-town skip"
    Then it runs the commands
      | BRANCH | COMMAND                 |
      | beta   | git rebase --abort      |
      |        | git checkout gamma      |
      | gamma  | git rebase origin/gamma |
      |        | git checkout main       |
      | main   | git push --tags         |
      |        | git stash pop           |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE            |
      | main   | local, remote | main commit        |
      | alpha  | local, remote | alpha commit       |
      | beta   | local         | local beta commit  |
      |        | remote        | remote beta commit |
      | gamma  | local, remote | gamma commit       |

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And my uncommitted file is stashed
    And my repo still has a rebase in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND                 |
      | beta   | git rebase --continue   |
      |        | git push                |
      |        | git checkout gamma      |
      | gamma  | git rebase origin/gamma |
      |        | git checkout main       |
      | main   | git push --tags         |
      |        | git stash pop           |
    And all branches are now synchronized
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And there is no rebase in progress anymore

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH | COMMAND                 |
      | beta   | git push                |
      |        | git checkout gamma      |
      | gamma  | git rebase origin/gamma |
      |        | git checkout main       |
      | main   | git push --tags         |
      |        | git stash pop           |
