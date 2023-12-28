Feature: add Git Town aliases to existing Git aliases

  Scenario: existing alias for "git append"
    Given I ran "git config --global alias.append checkout"
    When I run "git-town aliases add"
    Then it runs the commands
      | COMMAND                                                      |
      | git config --global alias.diff-parent "town diff-parent"     |
      | git config --global alias.hack "town hack"                   |
      | git config --global alias.kill "town kill"                   |
      | git config --global alias.prepend "town prepend"             |
      | git config --global alias.propose "town propose"             |
      | git config --global alias.rename-branch "town rename-branch" |
      | git config --global alias.repo "town repo"                   |
      | git config --global alias.ship "town ship"                   |
      | git config --global alias.sync "town sync"                   |
    And global Git setting "alias.append" is still "checkout"
    And global Git setting "alias.diff-parent" is now "town diff-parent"
    And global Git setting "alias.hack" is now "town hack"
    And global Git setting "alias.kill" is now "town kill"
    And global Git setting "alias.prepend" is now "town prepend"
    And global Git setting "alias.propose" is now "town propose"
    And global Git setting "alias.rename-branch" is now "town rename-branch"
    And global Git setting "alias.repo" is now "town repo"
    And global Git setting "alias.ship" is now "town ship"
    And global Git setting "alias.sync" is now "town sync"

    When I run "git-town aliases remove"
    Then it runs the commands
      | COMMAND                                         |
      | git config --global --unset alias.diff-parent   |
      | git config --global --unset alias.hack          |
      | git config --global --unset alias.kill          |
      | git config --global --unset alias.prepend       |
      | git config --global --unset alias.propose       |
      | git config --global --unset alias.rename-branch |
      | git config --global --unset alias.repo          |
      | git config --global --unset alias.ship          |
      | git config --global --unset alias.sync          |
    And global Git setting "alias.append" is still "checkout"
    And global Git setting "alias.diff-parent" is now ""
    And global Git setting "alias.hack" is still ""
    And global Git setting "alias.kill" is now ""
    And global Git setting "alias.prepend" is now ""
    And global Git setting "alias.propose" is now ""
    And global Git setting "alias.rename-branch" is now ""
    And global Git setting "alias.repo" is now ""
    And global Git setting "alias.ship" is now ""
    And global Git setting "alias.sync" is now ""
