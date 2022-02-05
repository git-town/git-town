@skipWindows
Feature: warn about unfinished prompt asking the user how to proceed

  Background:
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main   | local    | conflicting local commit  | conflicting_file | conflicting local content  |
      |        | remote   | conflicting remote commit | conflicting_file | conflicting remote content |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    And I run "git-town sync"
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """

  Scenario: attempt to sync again and choosing to quit
    When I run "git-town sync" and answer the prompts:
      | PROMPT                       | ANSWER  |
      | Please choose how to proceed | [ENTER] |
    Then it runs no commands
    And it prints:
      """
      You have an unfinished `sync` command that ended on the `main` branch now.
      """
    And my uncommitted file is stashed

  Scenario: attempt to sync again and choosing to continue without resolving conflicts
    When I run "git-town sync" and answer the prompts:
      | PROMPT                       | ANSWER        |
      | Please choose how to proceed | [DOWN][ENTER] |
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And my uncommitted file is stashed

  Scenario: attempt to sync again and choosing to continue after resolving conflicts
    When I resolve the conflict in "conflicting_file"
    And I run "git-town sync", answer the prompts, and close the next editor:
      | PROMPT                       | ANSWER        |
      | Please choose how to proceed | [DOWN][ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git rebase --continue              |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |
    And all branches are now synchronized

  Scenario: attempt to sync again and choosing to abort
    When I run "git-town sync" and answer the prompts:
      | PROMPT                       | ANSWER              |
      | Please choose how to proceed | [DOWN][DOWN][ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | main    | git rebase --abort   |
      |         | git checkout feature |
      | feature | git stash pop        |
    And my repo is left with my original commits

  Scenario: run another command after manually aborting
    When I run "git rebase --abort"
    And I run "git checkout feature"
    And I run "git stash pop"
    And I run "git-town kill" and answer the prompts:
      | PROMPT                       | ANSWER                    |
      | Please choose how to proceed | [DOWN][DOWN][DOWN][ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                        |
      | feature | git fetch --prune --tags       |
      |         | git push origin :feature       |
      |         | git add -A                     |
      |         | git commit -m "WIP on feature" |
      |         | git checkout main              |
      | main    | git branch -D feature          |

  Scenario: does not report unfinished state after abort
    When I run "git-town abort"
    And I run "git-town kill"
    Then it does not print "You have an unfinished `sync` command that ended on the `main` branch now."
