Feature: renaming a perennial branch with a tracking branch

  Background:
    Given my repo has the perennial branches "qa" and "production"
    And my repo has a feature branch named "child-feature" as a child of "production"
    And the following commits exist in my repo
      | BRANCH        | LOCATION      | MESSAGE              |
      | production    | local, remote | production commit    |
      | child-feature | local, remote | child feature commit |
    And I am on the "production" branch

  Scenario: normal rename fails
    When I run "git-town rename-branch production renamed-production"
    Then it runs no commands
    And it prints the error:
      """
      "production" is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'
      """

  Scenario: forced rename works
    When I run "git-town rename-branch --force production renamed-production"
    Then it runs the commands
      | BRANCH             | COMMAND                                  |
      | production         | git fetch --prune --tags                 |
      |                    | git branch renamed-production production |
      |                    | git checkout renamed-production          |
      | renamed-production | git push -u origin renamed-production    |
      |                    | git push origin :production              |
      |                    | git branch -D production                 |
    And I am now on the "renamed-production" branch
    And the perennial branches are now configured as "qa" and "renamed-production"
    And my repo now has the following commits
      | BRANCH             | LOCATION      | MESSAGE              |
      | child-feature      | local, remote | child feature commit |
      | renamed-production | local, remote | production commit    |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT             |
      | child-feature | renamed-production |

  Scenario: undo
    Given I run "git-town rename-branch --force production renamed-production"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH             | COMMAND                                             |
      | renamed-production | git branch production {{ sha 'production commit' }} |
      |                    | git push -u origin production                       |
      |                    | git push origin :renamed-production                 |
      |                    | git checkout production                             |
      | production         | git branch -D renamed-production                    |
    And I am now on the "production" branch
    And the perennial branches are now configured as "qa" and "production"
    And my repo now has the following commits
      | BRANCH        | LOCATION      | MESSAGE              |
      | production    | local, remote | production commit    |
      | child-feature | local, remote | child feature commit |
