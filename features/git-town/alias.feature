Feature: git town: alias

  Scenario: add alias
    When I run `gt alias true`
    Then it runs the commands
      | COMMAND                                                           |
      | git config --global alias.append "!gt append"                     |
      | git config --global alias.hack "!gt hack"                         |
      | git config --global alias.kill "!gt kill"                         |
      | git config --global alias.new-pull-request "!gt new-pull-request" |
      | git config --global alias.prepend "!gt prepend"                   |
      | git config --global alias.prune-branches "!gt prune-branches"     |
      | git config --global alias.rename-branch "!gt rename-branch"       |
      | git config --global alias.repo "!gt repo"                         |
      | git config --global alias.ship "!gt ship"                         |
      | git config --global alias.sync "!gt sync"                         |


  Scenario: remove alias
    Given I run `gt alias true`
    When I run `gt alias false`
    Then it runs the commands
      | COMMAND                                            |
      | git config --global --unset alias.append           |
      | git config --global --unset alias.hack             |
      | git config --global --unset alias.kill             |
      | git config --global --unset alias.new-pull-request |
      | git config --global --unset alias.prepend          |
      | git config --global --unset alias.prune-branches   |
      | git config --global --unset alias.rename-branch    |
      | git config --global --unset alias.repo             |
      | git config --global --unset alias.ship             |
      | git config --global --unset alias.sync             |


  Scenario: remove alias does not remove unrelated aliases
    Given I run `git config --global alias.hack checkout`
    When I run `gt alias false`
    Then it runs no commands


  Scenario: invalid value
    When I run `gt alias other`
    Then I get the error "Invalid value: 'other'"
    And I get the error
      """
      Usage:
        gt alias (true | false) [flags]
      """
