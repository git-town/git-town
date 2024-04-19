Feature: ship a parent branch

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And the current branch is "child"
    When I run "git-town ship parent -m 'parent done'"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                        |
      | child  | git fetch --prune --tags       |
      |        | git checkout main              |
      | main   | git merge --squash --ff parent |
      |        | git commit -m "parent done"    |
      |        | git push                       |
      |        | git branch -D parent           |
      |        | git checkout child             |
    And the current branch is now "child"
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
      | child  | git checkout main                           |
      | main   | git revert {{ sha 'parent done' }}          |
      |        | git push                                    |
      |        | git branch parent {{ sha 'parent commit' }} |
      |        | git checkout child                          |
    And the current branch is now "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | parent done          |
      |        |               | Revert "parent done" |
      | child  | local, origin | child commit         |
      | parent | local, origin | parent commit        |
    And the initial branches and lineage exist
