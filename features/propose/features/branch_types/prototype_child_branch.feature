@skipWindows
Feature: Create proposals for prototype branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS |
      | parent    | feature   | main   | local     |
      | prototype | prototype | parent | local     |
    And the current branch is "prototype"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    And a proposal for this branch does not exist
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                                                                       |
      | prototype | git fetch --prune --tags                                                      |
      | <none>    | Looking for proposal online ... ok                                            |
      | prototype | git checkout main                                                             |
      | main      | git rebase origin/main --no-update-refs                                       |
      |           | git checkout parent                                                           |
      | parent    | git merge --no-edit --ff main                                                 |
      |           | git push -u origin parent                                                     |
      |           | git checkout prototype                                                        |
      | prototype | git merge --no-edit --ff parent                                               |
      | <none>    | open https://github.com/git-town/git-town/compare/parent...prototype?expand=1 |
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parent...prototype?expand=1
      """
    And it prints:
      """
      branch "prototype" is no longer a prototype branch
      """
    And there are now no prototype branches
