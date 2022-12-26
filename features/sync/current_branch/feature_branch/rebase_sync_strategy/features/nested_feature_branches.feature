Feature: nested feature branches

  Background:
    Given setting "sync-strategy" is "rebase"
    And a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | child  | git fetch --prune --tags    |
      |        | git checkout main           |
      | main   | git rebase origin/main      |
      |        | git push                    |
      |        | git checkout parent         |
      | parent | git rebase origin/parent    |
      |        | git rebase main             |
      |        | git push --force-with-lease |
      |        | git checkout child          |
      | child  | git rebase origin/child     |
      |        | git rebase parent           |
      |        | git push --force-with-lease |
    And all branches are now synchronized
    And the current branch is still "child"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | child  | local, origin | origin main commit   |
      |        |               | local main commit    |
      |        |               | origin parent commit |
      |        |               | local parent commit  |
      |        |               | origin child commit  |
      |        |               | local child commit   |
      | parent | local, origin | origin main commit   |
      |        |               | local main commit    |
      |        |               | origin parent commit |
      |        |               | local parent commit  |
