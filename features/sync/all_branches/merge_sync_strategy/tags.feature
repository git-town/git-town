Feature: sync all branches syncs the tags

  Scenario:
    Given a Git repo clone
    Given the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "main"
    When I run "git-town sync --all"
    Then these tags exist
      | NAME       | LOCATION      |
      | local-tag  | local, origin |
      | origin-tag | local, origin |
