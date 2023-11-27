Feature: shorten Git Town commands

  Scenario: inside a Git repo
    When I run "git-town aliases add"
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

  Scenario: outside a Git repo
    Given I am outside a Git repo
    When I run "git-town aliases add"
    Then it does not print "not a git repository"
