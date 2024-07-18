Feature: sync the current observed branch in a local repo

  Background:
    Given a local Git repo clone
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS |
      | prototype | prototype | main   | local     |
    And the commits
      | BRANCH    | LOCATION | MESSAGE      | FILE NAME  |
      | main      | local    | main commit  | main_file  |
      | prototype | local    | local commit | local_file |
    And the current branch is "prototype"
    And an uncommitted file
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                       |
      | prototype | git add -A                    |
      |           | git stash                     |
      |           | git merge --no-edit --ff main |
      |           | git stash pop                 |
    And these commits exist now
      | BRANCH    | LOCATION | MESSAGE                            |
      | main      | local    | main commit                        |
      | prototype | local    | local commit                       |
      |           |          | main commit                        |
      |           |          | Merge branch 'main' into prototype |
    And all branches are now synchronized
    And the current branch is still "prototype"
    And the uncommitted file still exists
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH    | COMMAND       |
      | prototype | git add -A    |
      |           | git stash     |
      |           | git stash pop |
    And the current branch is still "prototype"
    And the initial commits exist
    And the initial branches and lineage exist
