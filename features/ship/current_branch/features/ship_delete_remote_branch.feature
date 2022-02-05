Feature: ship-delete-remote-branch disabled

  Background:
    Given my repo has a local feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, remote | feature commit |
    And I am on the "feature" branch
    And my repo has "git-town.ship-delete-remote-branch" set to "false"
    When I run "git-town ship -m 'feature done'"
    And the remote deletes the "feature" branch

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git branch -D feature              |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES |
      | local, remote | main     |
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, remote | feature done |
    And Git Town now knows no branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch feature {{ sha 'feature commit' }} |
      |         | git revert {{ sha 'feature done' }}           |
      |         | git push                                      |
      |         | git checkout feature                          |
      | feature | git checkout main                             |
      | main    | git checkout feature                          |
    And I am now on the "feature" branch
    And my repo now has the commits
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, remote | feature done          |
      |         |               | Revert "feature done" |
      | feature | local         | feature commit        |
    And my repo now has its initial branches and branch hierarchy
