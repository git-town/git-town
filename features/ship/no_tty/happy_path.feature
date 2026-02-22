@skipWindows
Feature: no TTY, no main branch

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | existing | local    | existing commit |
    And Git setting "git-town.ship-strategy" is "fast-forward"
    And the current branch is "existing"
    When I run "git-town ship" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                      |
      | existing | git checkout main            |
      | main     | git merge --ff-only existing |
      |          | git branch -D existing       |

  Scenario: undo
    When I run "git-town undo" in a non-TTY shell
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git branch existing {{ sha 'existing commit' }} |
      |        | git checkout existing                           |
