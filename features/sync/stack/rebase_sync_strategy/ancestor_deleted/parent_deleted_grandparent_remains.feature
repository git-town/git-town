Feature: the parent of a branch was deleted remotely but a grandparent remains

  Background:
    Given a Git repo with origin
    And the branches
      | NAME        | TYPE    | PARENT      | LOCATIONS     |
      | grandparent | feature | main        | local, origin |
      | parent      | feature | grandparent | local, origin |
      | child       | feature | parent      | local, origin |
    And the commits
      | BRANCH      | LOCATION      | MESSAGE            |
      | grandparent | local, origin | grandparent commit |
      | parent      | local, origin | parent commit      |
      | child       | local, origin | child commit       |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin deletes the "parent" branch
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH      | COMMAND                                                         |
      | child       | git fetch --prune --tags                                        |
      |             | git checkout grandparent                                        |
      | grandparent | git checkout child                                              |
      | child       | git pull                                                        |
      |             | git -c rebase.updateRefs=false rebase --onto grandparent parent |
      |             | git push --force-with-lease                                     |
      |             | git branch -D parent                                            |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        grandparent
          child
      """
    And these commits exist now
      | BRANCH      | LOCATION      | MESSAGE            |
      | grandparent | local, origin | grandparent commit |
      | child       | local, origin | child commit       |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git reset --hard {{ sha 'child commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git branch parent {{ sha 'parent commit' }}     |
    Then the initial branches and lineage exist now
    And the initial commits exist now
