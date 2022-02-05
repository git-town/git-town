Feature: on a forked repo

  Background:
    Given my repo has an upstream repo
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/main   |
      |        | git fetch upstream main  |
      |        | git rebase upstream/main |
      |        | git push                 |
      |        | git branch new main      |
      |        | git checkout new         |
      | new    | git stash pop            |
    And I am now on the "new" branch
    And my workspace still contains my uncommitted file
    And my repo now has the commits
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, remote, upstream | upstream commit |
      | new    | local                   | upstream commit |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | new    | git add -A        |
      |        | git stash         |
      |        | git checkout main |
      | main   | git branch -d new |
      |        | git stash pop     |
    And I am now on the "main" branch
    And my repo now has the commits
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, remote, upstream | upstream commit |
    And Git Town now knows about no branch hierarchy
