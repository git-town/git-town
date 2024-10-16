Feature: sync perennial branch that was deleted at the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME        | TYPE      | PARENT      | LOCATIONS     |
      | perennial-1 | perennial |             | local, origin |
      | perennial-2 | perennial |             | local, origin |
      | feature-1a  | feature   | perennial-1 | local, origin |
      | feature-1b  | feature   | perennial-1 | local, origin |
      | feature-2a  | feature   | perennial-2 | local, origin |
    And origin deletes the "perennial-1" branch
    And the current branch is "perennial-1"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                   |
      | perennial-1 | git fetch --prune --tags  |
      |             | git checkout main         |
      | main        | git branch -D perennial-1 |
      |             | git push --tags           |
    And it prints:
      """
      deleted branch "perennial-1"
      """
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES                                              |
      | local, origin | main, feature-1a, feature-1b, feature-2a, perennial-2 |
    And the perennial branches are now "perennial-2"
    And this lineage exists now
      | BRANCH     | PARENT      |
      | feature-2a | perennial-2 |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git branch perennial-1 {{ sha 'initial commit' }} |
      |        | git checkout perennial-1                          |
    And the current branch is now "perennial-1"
    And the initial branches and lineage exist now
    And the perennial branches are now "perennial-1" and "perennial-2"
