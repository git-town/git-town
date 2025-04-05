Feature: sync all branches with an unpushed tag and enabled push hook

  Background:
    Given a Git repo with origin
    And the tags
      | NAME      | LOCATION |
      | local-tag | local    |
    And the current branch is "main"
    And Git setting "git-town.push-hook" is "true"
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git push --tags          |
