Feature: proposing a stack containing a parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | parked  | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And the current branch is "child"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                   |
      | child  | git fetch --prune --tags                                                  |
      |        | git merge --no-edit --ff parent                                           |
      |        | git push                                                                  |
      | (none) | open https://github.com/git-town/git-town/compare/parent?expand=1         |
      |        | open https://github.com/git-town/git-town/compare/parent...child?expand=1 |
    And the initial branches exist now
    And the initial lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git reset --hard {{ sha 'child commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
    And the initial branches exist now
    And the initial lineage exists now
