Feature: local branch

  Background:
    Given a local Git repo
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | dead  | feature | main   | local     |
      | other | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | dead   | local    | dead commit  |
      | other  | local    | other commit |
    And the current branch is "dead"
    And an uncommitted file
    When I run "git-town delete dead"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                   |
      | dead   | git add -A                                                |
      |        | git commit -m "Committing open changes on deleted branch" |
      |        | git checkout main                                         |
      | main   | git branch -D dead                                        |
    And the current branch is now "main"
    And no uncommitted files exist now
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
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                               |
      | main   | git branch dead {{ sha 'Committing open changes on deleted branch' }} |
      |        | git checkout dead                                                     |
      | dead   | git reset --soft HEAD~1                                               |
    And the current branch is now "dead"
    And the uncommitted file still exists
    And the initial commits exist now
    And the initial branches and lineage exist now
