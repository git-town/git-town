Feature: ship the current feature branch from a subfolder on the shipped branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME               |
      | feature | local, origin | feature commit | new_folder/feature_file |
    And the current branch is "feature"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git fetch --prune --tags    |
      |         | git checkout main           |
      | main    | git merge --ff-only feature |
      |         | git push                    |
      |         | git push origin :feature    |
      |         | git branch -D feature       |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'feature commit' }} |
      |        | git push -u origin feature                    |
      |        | git checkout feature                          |
    And the current branch is now "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | main    | local, origin | feature commit |
      | feature | local, origin | feature commit |
    And the initial branches and lineage exist now
