Feature: delete another than the current branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | good | feature | main   | local, origin |
      | dead | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT |
      | main   | local, origin | conflicting commit | conflicting_file | main content |
      | dead   | local, origin | dead-end commit    | file             | dead content |
      | good   | local, origin | good commit        | file             | good content |
    And the current branch is "good"
    And an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town delete dead"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | good   | git fetch --prune --tags    |
      |        | git add -A                  |
      |        | git stash                   |
      |        | git push origin :dead       |
      |        | git rebase --onto main dead |
      |        | git branch -D dead          |
      |        | git stash pop               |
    And Git Town prints the error:
      """
      conflicts between your uncommmitted changes and the main branch
      """
    And the current branch is still "good"
    And the uncommitted file has content:
      """
      <<<<<<< Updated upstream
      main content
      =======
      conflicting content
      >>>>>>> Stashed changes
      """
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND        |
      | good   | git stash drop |
    Then the branches are now
      | REPOSITORY    | BRANCHES   |
      | local, origin | main, good |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            |
      | main   | local, origin | conflicting commit |
      | good   | local         | good commit        |
      |        | origin        | good commit        |
    And this lineage exists now
      | BRANCH | PARENT |
      | good   | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                              |
      | good   | git add -A                                           |
      |        | git commit -m "Committing open changes to undo them" |
      |        | git branch dead {{ sha 'dead-end commit' }}          |
      |        | git push -u origin dead                              |
      |        | git stash pop                                        |
    And the current branch is still "good"
    And the uncommitted file does not exist anymore
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                              |
      | main   | local, origin | conflicting commit                   |
      | dead   | local, origin | dead-end commit                      |
      | good   | local         | Committing open changes to undo them |
      |        | origin        | good commit                          |
    And the initial branches and lineage exist now
