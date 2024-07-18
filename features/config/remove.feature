Feature: reset the configuration

  Scenario: with configuration
    Given a Git repo clone
    And the branches
      | NAME    | TYPE      | PARENT | LOCATIONS |
      | feature | feature   | main   | local     |
      | qa      | perennial |        | local     |
      | staging | perennial |        | local     |
    And the main branch is "main"
    And the current branch is "feature"
    And global Git setting "alias.hack" is "town hack"
    And global Git setting "alias.sync" is "town sync"
    And global Git setting "alias.append" is "commit --amend"
    When I run "git-town config remove"
    Then it runs the commands
      | COMMAND                                |
      | git config --global --unset alias.hack |
      | git config --global --unset alias.sync |
    And Git Town is no longer configured
    And global Git setting "alias.append" is still "commit --amend"

  Scenario: no configuration
    Given a Git repo clone
    And Git Town is not configured
    When I run "git-town config remove"
    Then Git Town is no longer configured
