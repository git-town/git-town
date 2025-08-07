Feature: create a new top-level feature branch in a clean workspace using the "compress" sync strategy when there are new commits on main

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE            |
      | main    | origin        | new commit         |
      | feature | local, origin | already compressed |
    And the current branch is "feature"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git fetch --prune --tags                          |
      |         | git checkout main                                 |
      | main    | git -c rebase.updateRefs=false rebase origin/main |
      |         | git checkout -b new                               |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE            |
      | main    | local, origin | new commit         |
      | feature | local, origin | already compressed |
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |
      | new     | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                     |
      | new     | git checkout main                           |
      | main    | git reset --hard {{ sha 'initial commit' }} |
      |         | git checkout feature                        |
      | feature | git branch -D new                           |
    And the initial commits exist now
    And the initial branches and lineage exist now
