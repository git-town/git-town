Feature: ship the supplied feature branch in a local repo

  Background:
    Given the feature branches "feature" and "other"
    And my repo does not have an origin
    And the commits
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME        |
      | feature | local    | feature commit | conflicting_file |
    And the current branch is "other"
    And an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                      |
      | other  | git add -A                   |
      |        | git stash                    |
      |        | git checkout main            |
      | main   | git merge --squash feature   |
      |        | git commit -m "feature done" |
      |        | git branch -D feature        |
      |        | git checkout other           |
      | other  | git stash pop                |
    And the current branch is now "other"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | main   | local    | feature done |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git add -A                                    |
      |        | git stash                                     |
      |        | git checkout main                             |
      | main   | git reset --hard {{ sha 'initial commit' }}   |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout other                            |
      | other  | git stash pop                                 |
    And the current branch is now "other"
    And the initial commits exist
    And the initial branches and lineage exist
