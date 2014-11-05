Feature: git-ship
  the current branch


  Scenario: feature branch not ahead of main
    Given I am on a feature branch
    When I run `git ship -m 'feature done'` while allowing errors
    Then I get the error "The branch 'feature' has no commits to merge into 'main'."
    And I end up on the "feature" branch


  Scenario: on the main branch
    Given I am on the main branch
    When I run `git ship -m 'feature done'` while allowing errors
    Then I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "main" branch
    And there are no commits
    And there are no open changes


  Scenario: on non feature branch
    Given non-feature branch configuration "qa, production"
    And I am on the "production" branch
    When I run `git ship -m 'feature done'` while allowing errors
    Then I get the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "production" branch
    And there are no commits
    And there are no open changes


  Scenario: with uncommitted changes
    Given I am on a feature branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship -m 'feature done'` while allowing errors
    Then I get the error "You should not ship while having open files in Git"
    And I am still on the feature branch
    And there are no commits
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"

