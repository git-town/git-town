Feature: does not sync branches that exist only on remotes

  Background:
    Given my repo has a feature branch "mine"
    And the origin has a feature branch "other"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE         |
      | main   | origin        | main commit     |
      | mine   | local, origin | my commit       |
      | other  | origin        | coworker commit |
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
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE                       |
      | main   | local, origin | main commit                   |
      | mine   | local, origin | my commit                     |
      |        |               | main commit                   |
      |        |               | Merge branch 'main' into mine |
      | other  | origin        | coworker commit               |
