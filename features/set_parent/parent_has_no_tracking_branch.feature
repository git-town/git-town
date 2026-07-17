Feature: set the parent to a branch that has no tracking branch

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local         |
      | branch-2 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local         | commit 1 |
      | branch-2 | local, origin | commit 2 |
    And the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | TITLE             | BODY | URL                       |
      | 92 | branch-2      | main          | branch-2 proposal |      | https://example.com/pr/92 |
    And local Git setting "git-town.sync-feature-strategy" is "merge"
    And the current branch is "branch-2"
    When I run "git-town set-parent branch-1"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                              |
      |          | Finding proposal from branch-2 into main ... #92 (branch-2 proposal) |
      | branch-2 | git push -u origin branch-1                                          |
      |          | Updating target branch of proposal #92 to branch-1 ... ok            |
    And this lineage exists now
      """
      main
        branch-1
          branch-2
      """
    And the branches are now
      | REPOSITORY    | BRANCHES                 |
      | local, origin | main, branch-1, branch-2 |
    And the proposals are now
      """
      url: https://example.com/pr/92
      number: 92
      source: branch-2
      target: branch-1
      body:
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                               |
      |          | Updating target branch of proposal #92 to main ... ok |
      | branch-2 | git push origin :branch-1                             |
    And the initial branches and lineage exist now
    And the proposals are now
      """
      url: https://example.com/pr/92
      number: 92
      source: branch-2
      target: main
      body:
      """
