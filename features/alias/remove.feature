Feature: removing aliases

  Scenario: all existing aliases are Git Town shortcuts
    Given I ran "git-town aliases add"
    When I run "git-town aliases remove"
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

  Scenario: some aliases are not related to Git Town
    Given I ran "git config --global alias.hack checkout"
    When I run "git-town aliases remove"
    Then it runs no commands
