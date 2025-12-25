Feature: sync all branches syncs the tags

  Scenario:
    Given a Git repo with origin
    And Git setting "git-town.sync-tags" is "false"
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "main"
    When I run "git-town sync --all"
    Then the initial tags exist now
