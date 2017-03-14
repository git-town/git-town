Feature: gt sync --all: does not sync remote only branches

  Background:
    Given I have a feature branch named "my-feature"
    And my coworker has a feature branch named "co-feature"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE         | FILE NAME     |
      | main       | remote           | main commit     | main_file     |
      | my-feature | local and remote | my commit       | my_file       |
      | co-feature | remote           | coworker commit | coworker_file |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `gt sync --all`


  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                               |
      | main       | git fetch --prune                     |
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
    And I still have my uncommitted file
    And all branches are now synchronized
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE                             | FILE NAME     |
      | main       | local and remote | main commit                         | main_file     |
      | my-feature | local and remote | my commit                           | my_file       |
      |            |                  | main commit                         | main_file     |
      |            |                  | Merge branch 'main' into my-feature |               |
      | co-feature | remote           | coworker commit                     | coworker_file |
