Feature: listing the configuration

  Scenario: everything is configured
    Given the main branch is configured as "main"
    And the perennial branches are configured as "qa" and "staging"
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        main

      Perennial branches:
        qa
        staging
      """

  Scenario: everything is configured and there are nested branches
    Given the main branch is configured as "main"
    And my repo has the perennial branches "qa" and "staging"
    And my repo has the feature branches "parent-feature" and "stand-alone-feature"
    And my repo has a feature branch named "child-feature" as a child of "parent-feature"
    And my repo has a feature branch named "qa-hotfix" as a child of "qa"
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        main

      Perennial branches:
        qa
        staging

      Branch Ancestry:
        main
          parent-feature
            child-feature
          stand-alone-feature
        
        qa
          qa-hotfix
      """

  Scenario: the main branch is configured but the perennial branches are not
    Given the main branch is configured as "main"
    And the perennial branches are not configured
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        main

      Perennial branches:
        [none]
      """

  Scenario: the main branch is not configured but the perennial branches are
    Given the main branch name is not configured
    And the perennial branches are configured as "qa" and "staging"
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        [none]

      Perennial branches:
        qa
        staging
      """

  Scenario: nothing is configured yet
    Given I haven't configured Git Town yet
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        [none]

      Perennial branches:
        [none]
      """

  Scenario: disconnected branches
    Given the main branch is configured as "main"
    And my repo has a branch "parent-feature"
    And my repo has a feature branch named "child-feature" as a child of "parent-feature"
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        main

      Perennial branches:
        [none]

      Branch Ancestry:
        parent-feature
          child-feature
      """
