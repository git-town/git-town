@skipWindows
Feature: dry-run proposing changes

  Scenario: proposing changes
    Given tool "open" is installed
    And the current branch is a feature branch "feature"
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose --dry-run"
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --tags                                           |
      |         | git checkout main                                                  |
      | main    | git rebase origin/main                                             |
      |         | git checkout feature                                               |
      | feature | git merge --no-edit origin/feature                                 |
      |         | git merge --no-edit --ff main                                      |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And the current branch is still "feature"
    And the initial branches and lineage exist
