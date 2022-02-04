Feature: does not sync branches that exist only on the remote

  Background:
    Given my repo has a feature branch "my-feature"
    And a coworker has a feature branch "co-feature"
    And my repo contains the commits
      | BRANCH     | LOCATION      | MESSAGE         |
      | main       | remote        | main commit     |
      | my-feature | local, remote | my commit       |
      | co-feature | remote        | coworker commit |
    And I am on the "main" branch
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                               |
      | main       | git fetch --prune --tags              |
      |            | git rebase origin/main                |
      |            | git checkout my-feature               |
      | my-feature | git merge --no-edit origin/my-feature |
      |            | git merge --no-edit main              |
      |            | git push                              |
      |            | git checkout main                     |
      | main       | git push --tags                       |
    And I am still on the "main" branch
    And all branches are now synchronized
    And my repo now has the commits
      | BRANCH     | LOCATION      | MESSAGE                             |
      | main       | local, remote | main commit                         |
      | co-feature | remote        | coworker commit                     |
      | my-feature | local, remote | my commit                           |
      |            |               | main commit                         |
      |            |               | Merge branch 'main' into my-feature |
