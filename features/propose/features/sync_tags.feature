Feature: don't sync tags while proposing

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "feature"
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    And tool "open" is installed
    And Git Town setting "sync-tags" is "false"
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --no-tags                                        |
      |         | git checkout main                                                  |
      | main    | git rebase origin/main                                             |
      |         | git checkout feature                                               |
      | feature | git merge --no-edit --ff origin/feature                            |
      |         | git merge --no-edit --ff main                                      |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the initial commits exist
    And the initial lineage exists
    And the initial tags exist now
