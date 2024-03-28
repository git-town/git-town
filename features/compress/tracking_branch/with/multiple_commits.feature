Feature: compress the commits on a feature branch

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  |
      | feature | local, origin | commit 1 |
      |         |               | commit 2 |
      |         |               | commit 3 |
    When I run "git-town compress" and enter "compressed commit" for the commit message

  @debug @this
  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git reset main                                  |
      |         | git add -A                                      |
      |         | git commit                                      |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE           |
      | feature | local, origin | compressed commit |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                      |
      | feature | git reset --hard {{ sha 'local feature commit' }}                            |
      |         | git push --force-with-lease origin {{ sha 'origin feature commit' }}:feature |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | origin main commit    |
      |         |               | local main commit     |
      | feature | local         | local feature commit  |
      |         | origin        | origin feature commit |
    And the initial branches and lineage exist
