Feature: don't sync tags while proposing

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "feature"
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    And a proposal for this branch does not exist
    And tool "open" is installed
    And Git Town setting "sync-tags" is "false"
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --no-tags                                        |
      | <none>  | Looking for proposal online ... ok                                 |
      | feature | git checkout main                                                  |
      | main    | git rebase origin/main --no-update-refs                            |
      |         | git checkout feature                                               |
      | feature | git merge --no-edit --ff main                                      |
      |         | git merge --no-edit --ff origin/feature                            |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
    And the initial lineage exists now
    And the initial tags exist now
