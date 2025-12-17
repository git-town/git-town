Feature: ignore uncommitted changes using Git metadata

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-ignore-uncommitted" is "true"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                     |
      | feature | git add -A                                                  |
      |         | git stash -m "Git Town WIP"                                 |
      |         | git checkout main                                           |
      | main    | git -c color.ui=always merge --squash --ff feature          |
      |         | git commit -m "feature commit" --trailer "Co-authored-by: " |
      |         | git push                                                    |
      |         | git push origin :feature                                    |
      |         | git branch -D feature                                       |
      |         | git stash pop                                               |
      |         | git restore --staged .                                      |
    And the current branch is now "main"
    And the uncommitted file still exists
