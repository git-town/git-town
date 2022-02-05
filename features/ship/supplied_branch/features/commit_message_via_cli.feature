Feature: provide the commit message via a CLI argument

  Background:
    Given my repo has the feature branches "feature" and "other"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME        |
      | feature | local, remote | feature commit | conflicting_file |
    And I am on the "other" branch
    And my workspace has an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | other   | git fetch --prune --tags           |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git push origin :feature           |
      |         | git branch -D feature              |
      |         | git checkout other                 |
      | other   | git stash pop                      |
    And I am now on the "other" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY    | BRANCHES    |
      | local, remote | main, other |
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, remote | feature done |
    And Git Town now knows about this branch hierarchy
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
      |         | git push -u origin feature                    |
      |         | git revert {{ sha 'feature done' }}           |
      |         | git push                                      |
      |         | git checkout feature                          |
      | feature | git checkout main                             |
      | main    | git checkout other                            |
      | other   | git stash pop                                 |
    And I am now on the "other" branch
    And my repo now has the commits
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, remote | feature done          |
      |         |               | Revert "feature done" |
      | feature | local, remote | feature commit        |
    And my repo now has its initial branches and branch hierarchy
