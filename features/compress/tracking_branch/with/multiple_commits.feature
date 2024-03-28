Feature: compress the commits on a feature branch

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         |               | commit 2 | file_2    | content 2    |
      |         |               | commit 3 | file_3    | content 3    |
    When I run "git-town compress" and enter "compressed commit" for the commit message

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
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"
    And file "file_3" still has content "content 3"

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git reset --hard {{ sha 'commit 3' }}           |
      |         | git push --force-with-lease --force-if-includes |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
