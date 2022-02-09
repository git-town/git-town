Feature: ship the supplied feature branch in a local repo

  Background:
    Given the feature branches "feature" and "other"
    And my repo does not have an origin
    And the commits
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME        |
      | feature | local    | feature commit | conflicting_file |
    And I am on the "other" branch
    And my workspace has an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | other   | git add -A                   |
      |         | git stash                    |
      |         | git checkout feature         |
      | feature | git merge --no-edit main     |
      |         | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git branch -D feature        |
      |         | git checkout other           |
      | other   | git stash pop                |
    And I am now on the "other" branch
    And my workspace still contains my uncommitted file
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
    And now these commits exist
      | BRANCH | LOCATION | MESSAGE      |
      | main   | local    | feature done |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | other   | git add -A                                    |
      |         | git stash                                     |
      |         | git checkout main                             |
      | main    | git branch feature {{ sha 'feature commit' }} |
      |         | git revert {{ sha 'feature done' }}           |
      |         | git checkout feature                          |
      | feature | git checkout other                            |
      | other   | git stash pop                                 |
    And I am now on the "other" branch
    And now these commits exist
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | feature done          |
      |         |          | Revert "feature done" |
      | feature | local    | feature commit        |
    And the initial branch setup and hierarchy exist
