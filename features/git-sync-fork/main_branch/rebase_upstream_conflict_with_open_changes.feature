Feature: git-sync-fork: handling rebase conflicts between main branch and its remote with open changes

  Background:
    Given my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         | FILE NAME        | FILE CONTENT     |
      | main   | local    | local commit    | conflicting_file | local content    |
      |        | upstream | upstream commit | conflicting_file | upstream content |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git sync-fork`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                  |
      | main   | git stash -u             |
      |        | git fetch upstream       |
      |        | git rebase upstream/main |
    And I get the error
      """
      To abort, run "git sync-fork --abort".
      To continue after you have resolved the conflicts, run "git sync-fork --continue".
      """
    And my repo has a rebase in progress
    And my uncommitted file is still stashed away


  Scenario: aborting
    When I run `git sync-fork --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git stash pop      |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And there is no rebase in progress
    And I am left with my original commits


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync-fork --continue`
    Then it runs the Git commands
      | BRANCH | COMMAND               |
      | main   | git rebase --continue |
      |        | git push              |
      |        | git stash pop         |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And I still have the following commits
      | BRANCH | LOCATION                    | MESSAGE         | FILE NAME        |
      | main   | local, remote, and upstream | upstream commit | conflicting_file |
      |        | local and remote            | local commit    | conflicting_file |
    And now I have the following committed files
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git sync-fork --continue`
    Then it runs the Git commands
      | BRANCH | COMMAND       |
      | main   | git push      |
      |        | git stash pop |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And I still have the following commits
      | BRANCH | LOCATION                    | MESSAGE         | FILE NAME        |
      | main   | local, remote, and upstream | upstream commit | conflicting_file |
      |        | local and remote            | local commit    | conflicting_file |
    And now I have the following committed files
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | resolved content |
