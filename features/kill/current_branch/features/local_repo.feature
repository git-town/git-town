Feature: in a local repo

  Background:
    Given a local Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
      | other   | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
      | other   | local    | other commit   |
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                          |
      | feature | git add -A                                       |
      |         | git commit -m "Committing WIP for git town undo" |
      |         | git checkout main                                |
      | main    | git branch -D feature                            |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | other  | local    | other commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                         |
      | main    | git branch feature {{ sha 'Committing WIP for git town undo' }} |
      |         | git checkout feature                                            |
      | feature | git reset --soft HEAD~1                                         |
    And the current branch is now "feature"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
