@skipWindows
Feature: ship a coworker's feature branch

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE            | AUTHOR                            |
      | feature | local    | developer commit 1 | developer <developer@example.com> |
      |         |          | developer commit 2 | developer <developer@example.com> |
      |         |          | coworker commit    | coworker <coworker@example.com>   |

  Scenario: choose myself as the author
    When I run "git-town ship -m 'feature done'" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please choose an author for the squash commit | [ENTER] |
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                            |
      | main   | local, origin | feature done | developer <developer@example.com> |
    And Git Town is now aware of no branch hierarchy

  Scenario: choose a coworker as the author
    When I run "git-town ship -m 'feature done'" and answer the prompts:
      | PROMPT                                        | ANSWER        |
      | Please choose an author for the squash commit | [DOWN][ENTER] |
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                          |
      | main   | local, origin | feature done | coworker <coworker@example.com> |
    And Git Town is now aware of no branch hierarchy

  Scenario:  undo
    Given I ran "git-town ship -m 'feature done'" and answered the prompts:
      | PROMPT                                        | ANSWER  |
      | Please choose an author for the squash commit | [ENTER] |
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                        |
      | main    | git branch feature {{ sha 'coworker commit' }} |
      |         | git push -u origin feature                     |
      |         | git revert {{ sha 'feature done' }}            |
      |         | git push                                       |
      |         | git checkout feature                           |
      | feature | git checkout main                              |
      | main    | git checkout feature                           |
    And the current branch is now "feature"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | feature done          |
      |         |               | Revert "feature done" |
      | feature | local, origin | developer commit 1    |
      |         |               | developer commit 2    |
      |         |               | coworker commit       |
    And the initial branches and hierarchy exist
