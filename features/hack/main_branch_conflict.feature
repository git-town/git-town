Feature: git town-hack: resolving conflicts between main branch and its tracking branch

  To rely on Git Town to finish the command or cleanly abort
  When there are conflicting commits on the local and remote main branch
  I want to be given the choice to resolve the conflicts or abort.

  Background:
    Given my repo has a feature branch named "existing-feature"
    And the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
      |        | remote   | conflicting remote commit | conflicting_file | remote content |
    And I am on the "existing-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH           | COMMAND                  |
      | existing-feature | git fetch --prune --tags |
      |                  | git add -A               |
      |                  | git stash                |
      |                  | git checkout main        |
      | main             | git rebase origin/main   |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And my repo now has a rebase in progress
    And my uncommitted file is stashed

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH           | COMMAND                       |
      | main             | git rebase --abort            |
      |                  | git checkout existing-feature |
      | existing-feature | git stash pop                 |
    And I am now on the "existing-feature" branch
    And my workspace has the uncommitted file again
    And there is no rebase in progress
    And my repo is left with my original commits
    # And Git Town now has no branch hierarchy information TODO

  Scenario: continuing without resolving the conflicts
    When I run "git-town continue"
    Then it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And my uncommitted file is stashed
    And my repo still has a rebase in progress

  Scenario: continuing after resolving the conflicts but not finishing the rebase
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | main     | git rebase --continue    |
      |          | git push                 |
      |          | git branch new-feature main |
      |          | git checkout new-feature    |
      | new-feature | git stash pop            |
    And I am now on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH   | LOCATION      | MESSAGE                   | FILE NAME        |
      | main     | local, remote | conflicting remote commit | conflicting_file |
      |          |               | conflicting local commit  | conflicting_file |
      | new-feature | local         | conflicting remote commit | conflicting_file |
      |          |               | conflicting local commit  | conflicting_file |
    # And Git Town is now aware of this branch hierarchy  TODO
    #   | BRANCH      | PARENT |
    #   | new-feature | main   |

  Scenario: continuing after resolving the conflicts and finishing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | main     | git push                 |
      |          | git branch new-feature main |
      |          | git checkout new-feature    |
      | new-feature | git stash pop            |
    And I am now on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH      | LOCATION      | MESSAGE                   | FILE NAME        |
      | main        | local, remote | conflicting remote commit | conflicting_file |
      |             |               | conflicting local commit  | conflicting_file |
      | new-feature   | local         | conflicting remote commit | conflicting_file |
      |             |               | conflicting local commit  | conflicting_file |
    # And Git Town is now aware of this branch hierarchy  TODO
    #   | BRANCH      | PARENT |
    #   | new-feature | main   |
