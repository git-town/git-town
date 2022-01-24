Feature: git-town sync --all: does not sync remote only branches

  Background:
    Given my repo has a feature branch named "my-feature"
    And my coworker has a feature branch named "co-feature"
    And the following commits exist in my repo
      | BRANCH     | LOCATION      | MESSAGE         | FILE NAME     |
      | main       | remote        | main commit     | main_file     |
      | my-feature | local, remote | my commit       | my_file       |
      | co-feature | remote        | coworker commit | coworker_file |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                               |
      | main       | git fetch --prune --tags              |
      |            | git add -A                            |
      |            | git stash                             |
      |            | git rebase origin/main                |
      |            | git checkout my-feature               |
      | my-feature | git merge --no-edit origin/my-feature |
      |            | git merge --no-edit main              |
      |            | git push                              |
      |            | git checkout main                     |
      | main       | git push --tags                       |
      |            | git stash pop                         |
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE                             | FILE NAME     |
      | main       | local, remote | main commit                         | main_file     |
      | co-feature | remote        | coworker commit                     | coworker_file |
      | my-feature | local, remote | my commit                           | my_file       |
      |            |               | main commit                         | main_file     |
      |            |               | Merge branch 'main' into my-feature |               |
