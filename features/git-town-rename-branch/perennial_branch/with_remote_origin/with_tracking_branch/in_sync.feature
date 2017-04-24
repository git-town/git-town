Feature: git town-rename-branch: renaming a perennial branch with a tracking branch

  As a developer with a poorly named perennial branch
  I want to be able to rename it safely in one easy step
  So that the names of my branches match what they implement, and I can manage them effectively.


  Background:
    Given I have perennial branches named "qa" and "production"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE           |
      | main       | local and remote | main commit       |
      | production | local and remote | production commit |
      | qa         | local and remote | qa commit         |
    And I am on the "production" branch
    And I have an uncommitted file


  Scenario: error when trying to rename
    When I run `git town-rename-branch production renamed-production`
    Then it runs no commands
    And I get the error "The branch 'production' is not a feature branch."
    And I get the error "Run 'git town-rename-branch production renamed-production -f' to force the rename, then reconfigure git-town on any other clones of this repo."


  Scenario: forcing rename
    When I run `git town-rename-branch production renamed-production -f`
    Then it runs the commands
      | BRANCH             | COMMAND                                       |
      | production         | git fetch --prune                             |
      |                    | git add -A                                    |
      |                    | git stash                                     |
      |                    | git checkout -b renamed-production production |
      | renamed-production | git push -u origin renamed-production         |
      |                    | git push origin :production                   |
      |                    | git branch -D production                      |
      |                    | git stash pop                                 |
    And I end up on the "renamed-production" branch
    And my repo is configured with perennial branches as "qa" and "renamed-production"
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH             | LOCATION         | MESSAGE           |
      | main               | local and remote | main commit       |
      | qa                 | local and remote | qa commit         |
      | renamed-production | local and remote | production commit |


  Scenario: undo
    Given I have run `git town-rename-branch production renamed-production -f`
    When I run `git town-rename-branch --undo`
    Then it runs the commands
        | BRANCH             | COMMAND                                              |
        | renamed-production | git add -A                                           |
        |                    | git stash                                            |
        |                    | git branch production <%= sha 'production commit' %> |
        |                    | git push -u origin production                        |
        |                    | git push origin :renamed-production                  |
        |                    | git checkout production                              |
        | production         | git branch -d renamed-production                     |
        |                    | git stash pop                                        |
    And I end up on the "production" branch
    And my repo is configured with perennial branches as "qa" and "production"
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE           |
      | main       | local and remote | main commit       |
      | production | local and remote | production commit |
      | qa         | local and remote | qa commit         |
