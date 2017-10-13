Feature: git town: alias

  Scenario: add alias
    When I run `git-town alias true`
    Then Git Town runs the commands
      | COMMAND                                                            |
      | git config --global alias.append "town append"                     |
      | git config --global alias.hack "town hack"                         |
      | git config --global alias.kill "town kill"                         |
      | git config --global alias.new-pull-request "town new-pull-request" |
      | git config --global alias.prepend "town prepend"                   |
      | git config --global alias.prune-branches "town prune-branches"     |
      | git config --global alias.rename-branch "town rename-branch"       |
      | git config --global alias.repo "town repo"                         |
      | git config --global alias.ship "town ship"                         |
      | git config --global alias.sync "town sync"                         |


  Scenario: remove alias
    Given I run `git-town alias true`
    When I run `git-town alias false`
    Then Git Town runs the commands
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
    When I run `git-town alias false`
    Then Git Town runs no commands


  Scenario: invalid value
    When I run `git-town alias other`
    Then Git Town prints the error "Invalid value: 'other'"
    And it prints the error:
      """
      Usage:
        git-town alias (true | false) [flags]
      """


  Scenario: works outside of a Git repository
    Given my workspace is currently not a Git repository
    When I run `git-town alias true`
    Then Git Town does not print "Not a git repository"
