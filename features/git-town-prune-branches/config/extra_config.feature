Feature: pruning configuration data

  When having used Git Town for a while
  I want to prune unnecessary Git Town configuration information
  So that my Git configuration is lean and relevant.

  - running "git town prune configuration" prunes the configuration data


  Scenario: Git config contains information about non-existing branches
    Given my repository has the feature branches "existing-feature"
    And my Git configuration has the entries:
      | KEY | VALUE |
      |
