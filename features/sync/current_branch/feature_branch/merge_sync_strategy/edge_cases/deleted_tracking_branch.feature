Feature: restores deleted tracking branch

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And origin deletes the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune --tags   |
      |         | git checkout main          |
      | main    | git rebase origin/main     |
      |         | git checkout feature       |
      | feature | git merge --no-edit main   |
      |         | git push -u origin feature |
    And all branches are now synchronized
    And the current branch is still "feature"
    And now the initial commits exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git push origin :feature |
      |         | git checkout main        |
      | main    | git checkout feature     |
    And the current branch is still "feature"
    And now these commits exist
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And the initial branches and hierarchy exist
