Feature: does not sync branches that exist only on the remote

  Background:
    Given my repo has a feature branch "mine"
    And a coworker has a feature branch "other"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE         |
      | main   | remote        | main commit     |
      | mine   | local, remote | my commit       |
      | other  | remote        | coworker commit |
    And I am on the "main" branch
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                         |
      | main   | git fetch --prune --tags        |
      |        | git rebase origin/main          |
      |        | git checkout mine               |
      | mine   | git merge --no-edit origin/mine |
      |        | git merge --no-edit main        |
      |        | git push                        |
      |        | git checkout main               |
      | main   | git push --tags                 |
    And I am still on the "main" branch
    And all branches are now synchronized
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE                       |
      | main   | local, remote | main commit                   |
      | mine   | local, remote | my commit                     |
      |        |               | main commit                   |
      |        |               | Merge branch 'main' into mine |
      | other  | remote        | coworker commit               |
