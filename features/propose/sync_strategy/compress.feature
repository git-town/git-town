@skipWindows
Feature: proposing using the "compress" sync strategy

  Scenario: proposing changes
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                 |
      | existing | local    | local existing commit 1 |
      | existing | local    | local existing commit 2 |
      | existing | origin   | remote existing commit  |
    And the current branch is "existing"
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    And wait 1 second to ensure new Git timestamps
    And a proposal for this branch does not exist
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                             |
      | existing | git fetch --prune --tags                                            |
      |          | Looking for proposal online ... ok                                  |
      |          | git merge --no-edit --ff origin/existing                            |
      |          | git reset --soft main                                               |
      |          | git commit -m "local existing commit 1"                             |
      |          | git push --force-with-lease                                         |
      |          | open https://github.com/git-town/git-town/compare/existing?expand=1 |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                 |
      | existing | local, origin | local existing commit 1 |
