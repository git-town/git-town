Feature: offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And offline mode is enabled
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
      | other   | local, origin | other commit   |
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town delete"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                          |
      | feature | git add -A                                       |
      |         | git commit -m "Committing WIP for git town undo" |
      |         | git checkout main                                |
      | main    | git branch -D feature                            |
    And the current branch is now "main"
    And no uncommitted files exist now
    And the branches are now
      | REPOSITORY | BRANCHES             |
      | local      | main, other          |
      | origin     | main, feature, other |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | origin        | feature commit |
      | other   | local, origin | other commit   |
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
    And the initial commits exist now
    And the initial branches and lineage exist now
