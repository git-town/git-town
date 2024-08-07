Feature: don't sync the tags

  Background:
    Given a Git repo with origin
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "main"
    When I run "git-town sync --all --no-tags"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | main   | git fetch --prune --no-tags |
      |        | git rebase origin/main      |
    And these tags exist
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
