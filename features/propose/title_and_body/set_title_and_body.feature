@skipWindows
Feature: Prepopulate title and body

  Background:
    Given a Git repo with origin
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And a proposal for this branch does not exist
    And tool "open" is installed

  Scenario: provide title and body via CLI
    When I run "git-town propose --title=my_title --body=my_body"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                        |
      | feature | git fetch --prune --tags                                                                       |
      |         | Looking for proposal online ... ok                                                             |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1&title=my_title&body=my_body |

  Scenario: provide title via CLI
    When I run "git-town propose --title=my_title"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                           |
      | feature | git fetch --prune --tags                                                          |
      |         | Looking for proposal online ... ok                                                |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1&title=my_title |

  Scenario: provide body via CLI
    When I run "git-town propose --body=my_body"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                         |
      | feature | git fetch --prune --tags                                                        |
      |         | Looking for proposal online ... ok                                              |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1&body=my_body |

  Scenario: provide title via CLI and body via file
    And an uncommitted file "body.txt" with content:
      """
      Proposal
      body
      text!
      """
    When I run "git-town propose --body-file=body.txt"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                           |
      | feature | git fetch --prune --tags                                                                          |
      |         | git add -A                                                                                        |
      |         | git stash -m "Git Town WIP"                                                                       |
      |         | Looking for proposal online ... ok                                                                |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1&body=Proposal%0Abody%0Atext%21 |
      |         | git stash pop                                                                                     |
      |         | git restore --staged .                                                                            |

  Scenario: non-existing body file
    When I run "git-town propose --body-file zonk.txt"
    Then Git Town prints the error:
      """
      Error: open zonk.txt: no such file or directory
      """

  Scenario: provide title via CLI and body via STDIN
    When I pipe the following text into "git-town propose --body-file -":
      """
      Proposal
      body
      text
      """
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                        |
      | feature | git fetch --prune --tags                                                                       |
      |         | Looking for proposal online ... ok                                                             |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1&body=Proposal%0Abody%0Atext |
