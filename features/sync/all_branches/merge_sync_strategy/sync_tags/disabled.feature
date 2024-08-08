Feature: sync all branches syncs the tags

  Scenario:
    Given a Git repo with origin
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "main"
    And Git Town setting "sync-tags" is "false"
    When I run "git-town sync --all"
    Then the initial tags exist now
