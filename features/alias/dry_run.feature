Feature: dry-run prints the commands but does not add or remove aliases

  @this
  Scenario: dry-run adding aliases
    When I run "git-town aliases add --dry-run"
    Then it runs the commands
      | COMMAND                                                      |
      | git config --global alias.append "town append"               |
      | git config --global alias.diff-parent "town diff-parent"     |
      | git config --global alias.hack "town hack"                   |
      | git config --global alias.kill "town kill"                   |
      | git config --global alias.prepend "town prepend"             |
      | git config --global alias.propose "town propose"             |
      | git config --global alias.rename-branch "town rename-branch" |
      | git config --global alias.repo "town repo"                   |
      | git config --global alias.ship "town ship"                   |
      | git config --global alias.sync "town sync"                   |
    And global Git setting "alias.append" is still ""
    And global Git setting "alias.diff-parent" is still ""
    And global Git setting "alias.hack" is still ""
    And global Git setting "alias.kill" is still ""
    And global Git setting "alias.prepend" is still ""
    And global Git setting "alias.propose" is still ""
    And global Git setting "alias.rename-branch" is still ""
    And global Git setting "alias.repo" is still ""
    And global Git setting "alias.ship" is still ""
    And global Git setting "alias.sync" is still ""

  Scenario: dry-run removing aliases
    Given I ran "git-town aliases add"
    When I run "git-town aliases remove --dry-run"
    Then it runs the commands
      | COMMAND                                         |
      | git config --global --unset alias.append        |
      | git config --global --unset alias.diff-parent   |
      | git config --global --unset alias.hack          |
      | git config --global --unset alias.kill          |
      | git config --global --unset alias.prepend       |
      | git config --global --unset alias.propose       |
      | git config --global --unset alias.rename-branch |
      | git config --global --unset alias.repo          |
      | git config --global --unset alias.ship          |
      | git config --global --unset alias.sync          |
    And global Git setting "alias.append" is still "town append"
    And global Git setting "alias.diff-parent" is still "town diff-parent"
    And global Git setting "alias.hack" is still "town hack"
    And global Git setting "alias.kill" is still "town kill"
    And global Git setting "alias.prepend" is still "town prepend"
    And global Git setting "alias.propose" is still "town propose"
    And global Git setting "alias.rename-branch" is still "town rename-branch"
    And global Git setting "alias.repo" is still "town repo"
    And global Git setting "alias.ship" is still "town ship"
    And global Git setting "alias.sync" is still "town sync"
