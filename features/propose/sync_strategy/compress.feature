@skipWindows
Feature: proposing using the "compress" sync strategy

  Scenario: proposing changes
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                 |
      | existing | local    | local existing commit 1 |
      | existing | local    | local existing commit 2 |
      | existing | origin   | remote existing commit  |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "existing"
    And a proposal for this branch does not exist
    And tool "open" is installed
    And wait 1 second to ensure new Git timestamps
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                             |
      | existing | git fetch --prune --tags                                            |
      |          | git merge --no-edit --ff origin/existing                            |
      |          | git reset --soft main --                                            |
      |          | git commit -m "local existing commit 1"                             |
      |          | git push --force-with-lease                                         |
      |          | Looking for proposal online ... ok                                  |
      |          | open https://github.com/git-town/git-town/compare/existing?expand=1 |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                 |
      | existing | local, origin | local existing commit 1 |
