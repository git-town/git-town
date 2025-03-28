Feature: detaching a branch that contains merge commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      | branch-2 | local, origin | commit 2 |
    And I ran "git merge branch-1 --no-edit"
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-3 | local, origin | commit 3 |
    And the current branch is "branch-2"
    When I run "git-town detach"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git fetch --prune --tags                        |
      |          | git rebase --onto main branch-1                 |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git pull                                        |
      |          | git rebase --onto branch-1 branch-2             |
      |          | git push --force-with-lease                     |
      |          | git checkout branch-2                           |
    And Git Town prints the error
      """
      cannot detach in the presence of merge commits, please compress and try again
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git reset --hard {{ sha 'commit 2b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'commit 3b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-4                           |
      | branch-4 | git reset --hard {{ sha 'commit 4b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-5                           |
      | branch-5 | git reset --hard {{ sha 'commit 5b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the current branch is still "branch-2"
    And the initial commits exist now
    And the initial lineage exists now
