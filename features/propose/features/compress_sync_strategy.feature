@skipWindows
Feature: proposing using the "compress" sync strategy

  Scenario: proposing changes
    Given a Git repo with origin
    And the branch
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                 |
      | existing | local    | local existing commit 1 |
      | existing | local    | local existing commit 2 |
      | existing | origin   | remote existing commit  |
    And the current branch is "existing"
    And Git Town setting "sync-feature-strategy" is "compress"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    And wait 1 second to ensure new Git timestamps
    Given a proposal for this branch does not exist
    When I run "git-town propose"
    Then it runs the commands
      | BRANCH   | COMMAND                                                             |
      | existing | git fetch --prune --tags                                            |
      | <none>   | looking for proposal online ... ok                                  |
      | existing | git checkout main                                                   |
      | main     | git rebase origin/main                                              |
      |          | git checkout existing                                               |
      | existing | git merge --no-edit --ff origin/existing                            |
      |          | git merge --no-edit --ff main                                       |
      |          | git reset --soft main                                               |
      |          | git commit -m "local existing commit 1"                             |
      |          | git push --force-with-lease                                         |
      | <none>   | open https://github.com/git-town/git-town/compare/existing?expand=1 |
    And the current branch is still "existing"
    And the initial branches and lineage exist
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                 |
      | existing | local, origin | local existing commit 1 |
