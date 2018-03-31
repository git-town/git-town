Feature: git-town prune

  - runs all subcommands at once


  Background:
    Given my repository has the feature branches "active-feature" and "deleted-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                |
      | active-feature  | local and remote | active-feature commit  |
      | deleted-feature | local and remote | deleted-feature commit |
    And the "deleted-feature" branch gets deleted on the remote
    And I am on the "deleted-feature" branch
    And Git Town is aware of this branch hierarchy
      | BRANCH         | PARENT  |
      | active-feature | main    |
      | other-feature  | feature |
    And my workspace has an uncommitted file
    When I run `git-town prune`


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                       |
      | deleted-feature | git fetch --prune             |
      |                 | git checkout main             |
      | main            | git branch -D deleted-feature |
    And I end up on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES             |
      | local      | main, active-feature |
      | remote     | main, active-feature |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT |
      | active-feature | main   |
