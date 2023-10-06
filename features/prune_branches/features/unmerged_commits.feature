Feature: prune a branch with unmerged commits whose tracking branch was deleted

  Background:
    Given the feature branches "other" and "dead-end"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | other    | local, origin | other commit    |
      | dead-end | local, origin | dead-end commit |
    And origin deletes the "dead-end" branch
    And the current branch is "dead-end"
    And an uncommitted file
    When I run "git-town prune-branches"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | dead-end | git fetch --prune --tags |
      |          | git add -A               |
      |          | git stash                |
      |          | git checkout main        |
      | main     | git rebase origin/main   |
      |          | git merge --no-edit main |
      |          | git checkout dead-end    |
      | dead-end | git stash pop            |
    And it prints:
      """
      Branch "dead-end" was deleted at the remote but the local branch contains unshipped changes.
      """
    And the current branch is still "dead-end"
    And the uncommitted file still exists
    And now these commits exist
      | BRANCH   | LOCATION      | MESSAGE         |
      | dead-end | local         | dead-end commit |
      | other    | local, origin | other commit    |
    And the initial branches and hierarchy exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND       |
      | dead-end | git add -A    |
      |          | git stash     |
      |          | git stash pop |
    And the current branch is now "dead-end"
    And the uncommitted file still exists
    And the initial branches and hierarchy exist
