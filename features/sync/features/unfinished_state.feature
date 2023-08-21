@skipWindows
Feature: handle previously unfinished Git Town commands

  Background: a Git Town command stops unfinished
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
      |        | origin   | conflicting origin commit | conflicting_file | origin content |
    And an uncommitted file
    And I run "git-town sync"
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """

  Scenario: quit a command that is blocked by a previously unfinished Git Town command
    When I run "git-town sync" and answer the prompts:
      | PROMPT                       | ANSWER  |
      | Please choose how to proceed | [ENTER] |
    Then it runs no commands
    And it prints:
      """
      You have an unfinished `sync` command that ended on the `main` branch now.
      """
    And the uncommitted file is stashed

  Scenario: continue a previously unfinished Git Town command without resolving the conflict
    When I run "git-town sync" and answer the prompts:
      | PROMPT                       | ANSWER        |
      | Please choose how to proceed | [DOWN][ENTER] |
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the uncommitted file is stashed

  Scenario: resolve and run the command again
    When I resolve the conflict in "conflicting_file"
    And I run "git-town diff-parent", answer the prompts, and close the next editor:
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
  # notice how it executes the steps for "git sync" and not the steps for "git diff-parent" here

  Scenario: run a command and abort the previously unfinished one
    When I run "git-town sync" and answer the prompts:
      | PROMPT                       | ANSWER              |
      | Please choose how to proceed | [DOWN][DOWN][ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | main    | git rebase --abort   |
      |         | git checkout feature |
      | feature | git stash pop        |
    And now the initial commits exist

  Scenario: run a command, abort the previously finished one, and run another command
    When I run "git-town abort"
    And I run "git-town diff-parent"
    Then it does not print "You have an unfinished `sync` command that ended on the `main` branch now."

  # TODO: after updating to a godog version > 0.9, group this and the next Scenario Outline into a Rule block
  # and merge the common setup steps into a local Background block.
  @this
  Scenario Outline: commands that require the user to resolve a previously unfinished Git Town command
    When I run "git rebase --abort"
    And I run "git checkout feature"
    And I run "git stash pop"
    And I run "git-town <COMMAND>" and answer the prompts:
      | PROMPT                       | ANSWER  |
      | Please choose how to proceed | [ENTER] |
    Then it prints:
      """
      You have an unfinished `sync` command that ended on the `main` branch
      """

    Examples:
      | COMMAND           |
      | append foo        |
      | diff-parent       |
      | hack foo          |
      | new-pull-request  |
      | prepend foo       |
      | prune-branches    |
      | rename-branch foo |
      | set-parent        |
      | ship              |
      | switch            |
      | sync              |

  Scenario Outline: commands that don't require the user to resolve a previously unfinished Git Town command
    When I run "git rebase --abort"
    And I run "git checkout feature"
    And I run "git stash pop"
    And I run "git-town <COMMAND>"
    Then it runs without error

    Examples:
      | COMMAND                     |
      | aliases add                 |
      | config                      |
      | config main-branch          |
      | config offline              |
      | config perennial-branches   |
      | config pull-branch-strategy |
      | config push-hook            |
      | config push-new-branches    |
      | config reset                |
      | config sync-strategy        |
      | kill                        |
      | status reset                |
      | status                      |
      | version                     |
