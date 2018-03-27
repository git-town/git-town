Feature: warn about unfinished prompt asking the user how to proceed

  As a developer running a command after not finishing the last one
  I want to be warned about it and presented with options
  So I can finish work I started and discard old state that is now irrelevant

  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main   | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      |        | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    And I run `git-town sync`
    And it prints the error:
      """
      To abort, run "git-town sync --abort".
      To continue after you have resolved the conflicts, run "git-town sync --continue".
      """


  Scenario: attempting to sync again and choosing to quit
    When I run `git-town sync` and answer the prompts:
      | PROMPT                         | ANSWER  |
      | How would you like to proceed: | [ENTER] |
    Then it runs no commands
    And it prints the error "You have an unfinished `sync` command that ended on the `current-feature` branch 1 second ago."


  Scenario: attempting to sync again and choosing to continue without resolving conflicts
    When I run `git-town sync` and answer the prompts:
      | PROMPT                         | ANSWER        |
      | How would you like to proceed: | [DOWN][ENTER] |
    Then it runs no commands
    And it prints the error "You must resolve the conflicts before continuing"


  Scenario: attempting to sync again and choosing to continue after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town sync` and answer the prompts:
      | PROMPT                         | ANSWER        |
      | How would you like to proceed: | [DOWN][ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git rebase --continue              |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |


  Scenario: attempting to sync again and choosing to abort
    When I run `git-town sync` and answer the prompts:
      | PROMPT                         | ANSWER              |
      | How would you like to proceed: | [DOWN][DOWN][ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | main    | git rebase --abort   |
      |         | git checkout feature |
      | feature | git stash pop        |


  Scenario: running another command after manually aborting
    Given I discard the open changes
    And I checkout the "feature" branch
    When I run `git-town kill` and answer the prompts:
      | PROMPT                         | ANSWER                    |
      | How would you like to proceed: | [DOWN][DOWN][DOWN][ENTER] |
    Then it runs no commands
      | BRANCH  | COMMAND                        |
      | feature | git fetch --prune              |
      |         | git push origin :feature       |
      |         | git add -A                     |
      |         | git commit -m "WIP on feature" |
      |         | git checkout main              |
      | main    | git branch -D feature          |
