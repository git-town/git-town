Feature: shorter Git Town commands

  Scenario: add aliases
    When I run "git-town aliases add"
    Then it runs the commands
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

  Scenario: remove aliases
    Given I ran "git-town aliases add"
    When I run "git-town aliases remove"
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

  Scenario: removing aliases does not remove unrelated aliases
    Given I ran "git config --global alias.hack checkout"
    When I run "git-town aliases remove"
    Then it runs no commands

  Scenario: works outside Git repositories
    Given I am outside a Git repo
    When I run "git-town aliases add"
    Then it does not print "Not a git repository"
