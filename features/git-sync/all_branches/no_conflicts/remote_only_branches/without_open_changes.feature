Feature: git sync --all: does not sync remote only branches without open changes

  Background:
    Given I have a feature branch named "my-feature"
    And my coworker has a feature branch named "co-feature"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE         | FILE NAME     |
      | main       | remote           | main commit     | main_file     |
      | my-feature | local and remote | my commit       | my_file       |
      | co-feature | remote           | coworker commit | coworker_file |
    And I am on the "main" branch
    When I run `git sync --all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH     | COMMAND                               |
      | main       | git fetch --prune                     |
      | main       | git rebase origin/main                |
      | main       | git checkout my-feature               |
      | my-feature | git merge --no-edit origin/my-feature |
      | my-feature | git merge --no-edit main              |
      | my-feature | git push                              |
      | my-feature | git checkout main                     |
    And I am still on the "main" branch
    And all branches are now synchronized
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE                             | FILE NAME     |
      | main       | local and remote | main commit                         | main_file     |
      | my-feature | local and remote | Merge branch 'main' into my-feature |               |
      |            |                  | main commit                         | main_file     |
      |            |                  | my commit                           | my_file       |
      | co-feature | remote           | coworker commit                     | coworker_file |
