Feature: on a forked repo

  Background:
    Given my repo has an upstream repo
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                     |
      | main        | git fetch --prune --tags    |
      |             | git add -A                  |
      |             | git stash                   |
      |             | git rebase origin/main      |
      |             | git fetch upstream main     |
      |             | git rebase upstream/main    |
      |             | git push                    |
      |             | git branch new-feature main |
      |             | git checkout new-feature    |
      | new-feature | git stash pop               |
    And I am now on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the commits
      | BRANCH      | LOCATION                | MESSAGE         |
      | main        | local, remote, upstream | upstream commit |
      | new-feature | local                   | upstream commit |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH      | COMMAND                   |
      | new-feature | git add -A                |
      |             | git stash                 |
      |             | git checkout main         |
      | main        | git branch -D new-feature |
      |             | git stash pop             |
    And I am now on the "main" branch
    And my repo now has the commits
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, remote, upstream | upstream commit |
    And Git Town now has no branch hierarchy information
