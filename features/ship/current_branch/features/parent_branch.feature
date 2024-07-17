Feature: ship a parent branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    Given the current branch is "parent"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    When I run "git-town ship -m 'parent done'"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                        |
      | parent | git fetch --prune --tags       |
      |        | git checkout main              |
      | main   | git merge --squash --ff parent |
      |        | git commit -m "parent done"    |
      |        | git push                       |
      |        | git branch -D parent           |
    And it prints:
      """
      branch "child" is now a child of "main"
      """
    And the current branch is now "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | parent done   |
      | child  | local, origin | child commit  |
      | parent | origin        | parent commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | main   | git revert {{ sha 'parent done' }}          |
      |        | git push                                    |
      |        | git branch parent {{ sha 'parent commit' }} |
      |        | git checkout parent                         |
    And the current branch is now "parent"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | parent done          |
      |        |               | Revert "parent done" |
      | child  | local, origin | child commit         |
      | parent | local, origin | parent commit        |
    And the initial branches and lineage exist
