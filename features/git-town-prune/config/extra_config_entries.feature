Feature: pruning configuration data

  When having used Git Town for a while
  I want to prune unnecessary Git Town configuration information
  So that my Git configuration is lean and relevant.

  - running "git-town prune configuration" removes data of branches
    that don't exist in the local workspace


  Background:
    Given my repository has the feature branches "existing-feature"
    And Git Town is aware of this branch hierarchy
      | BRANCH              | PARENT  |
      | existing-feature    | main    |
      | nonexisting-feature | feature |
    When I run `git-town prune config`


  Scenario: result
    Then Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT |
      | existing-feature | main   |


  Scenario: undo
    When I run `git-town prune config --undo`
    Then Git Town is now aware of this branch hierarchy
      | BRANCH              | PARENT  |
      | existing-feature    | main    |
      | nonexisting-feature | feature |
