Feature: sync doesn't support --all and --strategy at the same time

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | origin        | main commit  |
      | alpha  | local, origin | alpha commit |
    And the current branch is "alpha"
    When I run "git-town sync --all --strategy=rebase"

  Scenario: result
    Then Git Town prints the error:
      """
      sync doesn't support --all and --strategy flags at the same time
      """
