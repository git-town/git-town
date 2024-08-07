Feature: sync all branches syncs the tags

  Scenario:
    Given a Git repo with origin
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "main"
    When I run "git-town sync --all"
    Then these tags exist now
      | NAME       | LOCATION      |
      | local-tag  | local, origin |
      | origin-tag | local, origin |
