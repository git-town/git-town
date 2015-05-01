Feature: git rename-branch: renaming a non-feature branch with a tracking branch (with open changes)

  As a developer with a poorly named non-feature branch
  I want to be able to rename it safely in one easy step
  So that I can stay organized and remain productive


  Background:
    Given I have branches named "qa" and "production"
    And my non-feature branches are configured as "qa" and "production"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE           |
      | main       | local and remote | main commit       |
      | production | local and remote | production commit |
      | qa         | local and remote | qa commit         |
    And I am on the "production" branch
    And I have an uncommitted file


  Scenario: error when trying to rename
    When I run `git rename-branch production renamed-production`
    Then it runs no Git commands
    And I get the error "The branch 'production' is not a feature branch."
    And I get the error "Run 'git rename-branch production renamed-production -f' to force the rename, then reconfigure git-town on any other clones of this repo."


  Scenario: forcing rename
    When I run `git rename-branch production renamed-production -f`
    Then it runs the Git commands
      | BRANCH             | COMMAND                                       |
      | production         | git fetch --prune                             |
      |                    | git stash -u                                  |
      |                    | git checkout -b renamed-production production |
      | renamed-production | git push -u origin renamed-production         |
      |                    | git push origin :production                   |
      |                    | git branch -D production                      |
      |                    | git stash pop                                 |
    And I end up on the "renamed-production" branch
    And my non-feature branches are now configured as "qa" and "renamed-production"
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH             | LOCATION         | MESSAGE           |
      | main               | local and remote | main commit       |
      | qa                 | local and remote | qa commit         |
      | renamed-production | local and remote | production commit |
