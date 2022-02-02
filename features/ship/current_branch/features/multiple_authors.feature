@skipWindows
Feature: shipping a coworker's feature branch

  Background:
    Given my repo has a feature branch "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE         | AUTHOR                            |
      | feature | local    | feature commit1 | developer <developer@example.com> |
      |         |          | feature commit2 | developer <developer@example.com> |
      |         |          | feature commit3 | coworker <coworker@example.com>   |
    And I am on the "feature" branch

  Scenario: choosing myself as the author
    When I run "git-town ship -m 'feature done'" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please choose an author for the squash commit | [ENTER] |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                            |
      | main   | local, remote | feature done | developer <developer@example.com> |
    And Git Town now has no branch hierarchy information

  Scenario: choosing my coworker as the author
    When I run "git-town ship -m 'feature done'" and answer the prompts:
      | PROMPT                                        | ANSWER        |
      | Please choose an author for the squash commit | [DOWN][ENTER] |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                          |
      | main   | local, remote | feature done | coworker <coworker@example.com> |
    And Git Town now has no branch hierarchy information

  Scenario:  undo
    Given I run "git-town ship -m 'feature done'" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please choose an author for the squash commit | [ENTER] |
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                        |
      | main    | git branch feature {{ sha 'feature commit3' }} |
      |         | git push -u origin feature                     |
      |         | git revert {{ sha 'feature done' }}            |
      |         | git push                                       |
      |         | git checkout feature                           |
      | feature | git checkout main                              |
      | main    | git checkout feature                           |
    And I am now on the "feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, remote | feature done          |
      |         |               | Revert "feature done" |
      | feature | local, remote | feature commit1       |
      |         |               | feature commit2       |
      |         |               | feature commit3       |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
